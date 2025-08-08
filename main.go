package main

import (
	"devapi/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 允许跨域
	r.Use(gin.HandlerFunc(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, cache-control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}))
	// 注册路由
	r.LoadHTMLGlob("./static/*.html")
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./static/favicon.ico")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/netease", api.Netease)
	r.GET("/single", api.Single)
	r.GET("/lyric", api.Lyric)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", nil)
	})
	// 启动服务
	r.Run(":15967") // 监听并在 15967 端口启动服务
}
