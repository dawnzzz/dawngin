package dain

import (
	"net/http"
)

// HandlerFunc handler 函数类型
type HandlerFunc func(c *Context)

type Engine struct {
	// 路由器
	router *router
	// 继承 RouterGroup，把根也看作是一个分组
	*RouterGroup
	// 记录所有的路由分组
	groups []*RouterGroup
}

// RouterGroup 分组路由
type RouterGroup struct {
	prefix     string        // 当前分组的公共前缀
	parent     *RouterGroup  // 记录当前分组的上一层
	middleware []HandlerFunc // 记录中间件
	engine     *Engine       // 记录所属的 Engine
}

// New 返回一个 Engine 指针
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup) // 在 Engine 中保存新的分组
	return newGroup
}

// addRouter 实现路由注册功能
func (group *RouterGroup) addRouter(method, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	group.engine.router.addRouter(method, pattern, handler)
}

// Get 路由注册 GET 请求方式
func (group *RouterGroup) Get(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

// Post 路由注册 POST 请求方式
func (group *RouterGroup) Post(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
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
