package main

import (
	"DawnGin/dain"
	"net/http"
)

func main() {
	e := dain.New()

	e.Get("/", func(c *dain.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Dawn</h1>")
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

	e.Run(":9000")
}
