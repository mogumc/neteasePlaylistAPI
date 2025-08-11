package main

import (
	"devapi/api"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("release")

	r := gin.New()
	// 设置日志输出
	r.Use(gin.Logger(), gin.Recovery())

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
	fmt.Println("Init Any* HTML Glob")
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./static/favicon.ico")
	})
	fmt.Println("Route /favicon.ico registered")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	fmt.Println("Route / registered")

	r.GET("/netease", api.Netease)
	fmt.Println("Route /netease registered")

	r.GET("/single", api.Single)
	fmt.Println("Route /single registered")

	r.GET("/lyric", api.Lyric)
	fmt.Println("Route /lyric registered")

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", nil)
	})
	fmt.Println("Route Any* registered")

	// 启动服务
	fmt.Println("Server is running on http://localhost:15967")
	r.Run(":15967") // 监听并在 15967 端口启动服务
}
