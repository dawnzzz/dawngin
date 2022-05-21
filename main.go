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

	e.Post("/login", func(c *dain.Context) {
		c.JSON(http.StatusOK, dain.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	e.Run(":9000")
}
