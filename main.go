package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jiaojiajun/scoped_cache_demo/dependency_utils"
	"github.com/jiaojiajun/scoped_cache_demo/handler"
)

var userHandler = handler.GetSingletonNameHandler()

func main() {
	r := gin.Default()
	r.Use(AutoWire())
	r.GET("/name", GetFirstName)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func AutoWire() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dependency_utils.WireApp(ctx)
	}
}

func GetFirstName(ctx *gin.Context) {
	var name, err = userHandler.GetFirstName(ctx)
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": "error",
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": name,
		})
	}

}
