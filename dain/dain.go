package dain

import (
	"fmt"
	"net/http"
)

// HandlerFunc handler 函数类型
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	// 路由器
	router map[string]HandlerFunc
}

// New 返回一个 Engine 指针
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

// AddRouter 实现路由注册功能
func (e *Engine) AddRouter(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

// Get 路由注册 GET 请求方式
func (e *Engine) Get(pattern string, handler HandlerFunc) {
	e.AddRouter("GET", pattern, handler)
}

// Post 路由注册 POST 请求方式
func (e *Engine) Post(pattern string, handler HandlerFunc) {
	e.AddRouter("POST", pattern, handler)
}

// 实现 http.Handler 接口，自定义路由器
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND FOR PATH: %v", r.URL.Path)
	}
}

// Run 运行服务器
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
