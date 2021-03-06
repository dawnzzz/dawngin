package main

import (
	"DawnGin/dain"
	"net/http"
)

func main() {
	// 默认使用 Logger 和 Recovery 中间件
	e := dain.Default()

	// 加载静态文件
	e.Static("/static", "./static")

	// 加载模板
	e.LoadHTMLGlob("templates/*")

	e.Get("/index", func(c *dain.Context) {
		c.HTML(http.StatusOK, "index.tmpl", c.Path)
	})

	// 测试 Recovery 中间件
	e.Get("/panic", func(c *dain.Context) {
		array := []int{1, 2, 3}
		c.JSON(http.StatusOK, dain.H{
			"msg": array[100],
		})
	})

	e.Get("/hello", func(c *dain.Context) {
		c.String(http.StatusOK, "Hello World, URL path = %v", c.Path)
	})

	e.Get("/hello/:name", func(c *dain.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello, you are %v, URL path = %v", name, c.Path)
	})

	e.Get("/file/*filename", func(c *dain.Context) {
		filename := c.Param("filename")
		c.JSON(http.StatusOK, dain.H{
			"filename": filename,
			"msg":      "OK",
		})
	})

	e.Post("/login", func(c *dain.Context) {
		c.JSON(http.StatusOK, dain.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	// 分组路由
	v1 := e.Group("/v1")
	{
		v1.Get("/video/:name", func(c *dain.Context) {
			videoName := c.Param("name")
			c.String(http.StatusOK, "Hello, this is v1 group, video name = %v, path = %v", videoName, c.Path)
		})
	}

	e.Run(":9000")
}
