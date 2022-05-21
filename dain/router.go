package dain

import (
	"fmt"
	"log"
)

type router struct {
	handlers map[string]HandlerFunc
}

func NewRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// addRouter 实现路由注册功能
func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	log.Printf("Router %v - %v\n", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle 实现路由功能
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		fmt.Fprintf(c.Writer, "404 NOT FOUND FOR PATH: %v", c.Path)
	}
}
