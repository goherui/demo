package service

import (
	"demo/goods/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "请上传文件",
		})
	}
	upload := pkg.Upload(file)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "文件上传成功",
		"URL":  upload,
	})
}
