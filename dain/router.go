package dain

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // 保存 trie 树的根，key为 Method，value为树根
	handlers map[string]HandlerFunc // 保存 pattern 与 HandlerFunc 的映射关系
}

func NewRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析 pattern，返回对应的 parts（路由中的一部分）
// 如 pattern 为 /hello/world，那么对应的 parts 为 []{"hello", "world"}
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0, len(vs))
	for _, item := range vs {
		if item != "" {
			// 不为空，加入到 parts 中
			parts = append(parts, item)
			if item[0] == '*' {
				// 遇到通配符直接退出
				break
			}
		}
	}

	return parts
}

// addRouter 实现路由注册功能
func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	log.Printf("Router %v - %v\n", method, pattern)
	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	parts := parsePattern(pattern)
	r.roots[method].insert(pattern, parts, 0) // 插入到前缀树中
	r.handlers[key] = handler
}

// getRoute 根据请求的 method 和 path，找到对应的前缀树叶子节点和路由参数
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 查找的 parts
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		// 方法没有定义路由，直接返回
		return nil, nil
	}

	n := root.search(searchParts, 0) // 查找叶子节点
	if n != nil {
		// 可以找到
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return n, params
	}

	return nil, nil
}

// handle 实现路由功能
func (r *router) handle(c *Context) {
	// 在前缀树种查找路由，获取路由参数
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		// 匹配路由
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND FOR PATH: %v", c.Path)
		})
	}

	c.Next()
}
