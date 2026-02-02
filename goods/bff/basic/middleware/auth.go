package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token不能为空",
			})
			c.Abort()
			return
		}
		parseToken, s := ParseToken(token)
		if s != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token解析失败",
			})
			c.Abort()
			return
		}
		if parseToken == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token错误",
			})
			c.Abort()
			return
		}
		c.Set("userId", parseToken["userId"])
		c.Next()
	}
}
