package dain

import (
	"html/template"
	"net/http"
	"path"
	"strings"
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
	// 模板
	htmlTemplates *template.Template
	funcMap       template.FuncMap
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

// Default 使用 Logger 和 Recovery 中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())

	return engine
}

// Group 分组
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

// Use 为分组添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middleware = append(group.middleware, middlewares...)
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
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middleware...)
		}
	}
	c.engine = e
	e.router.handle(c)
}

// Run 运行服务器
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查文件是否可以访问
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 将磁盘上的某个路径 root 映射到 relativePath 上
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// 注册
	group.Get(urlPattern, handler)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
