package router

import (
	"demo/goods/bff/basic/middleware"
	"demo/goods/bff/handler/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	file := r.Group("/file")
	{
		file.POST("/upload", service.Upload)
	}
	user := r.Group("/user")
	{
		user.POST("/login", middleware.NegotiationCacheMiddleware("Etag"), service.Login)
		user.POST("/token", service.CreateToken)
	}
	goods := r.Group("/goods")
	{
		goods.POST("/create", service.GoodsCreate)
	}
	return r
}
