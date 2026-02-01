package service

import (
	"demo/goods/bff/basic/config"
	"demo/goods/bff/basic/middleware"
	__ "demo/goods/bff/basic/proto"
	"demo/goods/bff/handler/request"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var form request.Login
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}
	r, err := config.GoodsClient.Login(c, &__.LoginReq{
		Username: form.Username,
		Password: form.Password,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	userId, ok := r.UserMap["userId"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取用户 ID 失败",
			"code": http.StatusInternalServerError,
		})
		return
	}
	refreshToken, err := middleware.RefreshTokenHandler(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "生成刷新令牌失败",
			"code": http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  r.Msg,
		"code": r.Code,
		"data": gin.H{
			"userInfo":     r.UserMap,    // 用户信息
			"refreshToken": refreshToken, // 刷新令牌（7 天过期，用于换取新的 accessToken）
		},
	})
}
