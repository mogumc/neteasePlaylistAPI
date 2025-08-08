package main

import (
	"devapi/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 注册路由
	r.LoadHTMLGlob("./static/*.html")
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./static/favicon.ico")
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/netease", api.Netease)
	r.GET("/lyric", api.Lyric)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", nil)
	})
	// 启动服务
	r.Run(":15967") // 监听并在 15967 端口启动服务
}
