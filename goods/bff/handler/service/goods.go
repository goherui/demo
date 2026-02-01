package service

import (
	"demo/goods/bff/basic/config"
	__ "demo/goods/bff/basic/proto"
	"demo/goods/bff/handler/request"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GoodsCreate(c *gin.Context) {
	var form request.GoodsCreate
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}
	r, err := config.GoodsClient.GoodsCreate(c, &__.GoodsCreateReq{
		Title: form.Title,
		Price: float32(form.Price),
		Stock: int64(form.Stock),
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 将 Redis 数据添加到响应中
	c.JSON(http.StatusOK, gin.H{
		"error": r.Msg,
		"code":  r.Code,
		"goods": r.Goods,
	})
}
