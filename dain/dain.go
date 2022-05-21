package dain

import (
	"net/http"
)

// HandlerFunc handler 函数类型
type HandlerFunc func(c *Context)

type Engine struct {
	// 路由器
	router *router
}

// New 返回一个 Engine 指针
func New() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

// addRouter 实现路由注册功能
func (e *Engine) addRouter(method, pattern string, handler HandlerFunc) {
	e.router.addRouter(method, pattern, handler)
}

// Get 路由注册 GET 请求方式
func (e *Engine) Get(pattern string, handler HandlerFunc) {
	e.addRouter("GET", pattern, handler)
}

// Post 路由注册 POST 请求方式
func (e *Engine) Post(pattern string, handler HandlerFunc) {
	e.addRouter("POST", pattern, handler)
}

// 实现 http.Handler 接口，自定义路由器
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	e.router.handle(c)
}

// Run 运行服务器
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
