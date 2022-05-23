package dain

import "strings"

type node struct {
	pattern  string  // 待匹配的路由，只在叶子节点存储
	part     string  // 路由中的一部分
	children []*node // 叶子节点
	isWild   bool    // 若不是精确匹配则为 true；否则为 false
}

// matchChild 匹配第一个节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

// matchChildren 匹配所有的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}

	return children
}

// insert 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 到达叶子节点
		// 仅仅在叶子节点保存 pattern 信息
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// 若当前层没有匹配，新建一个节点，并插入到孩子节点中
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 向下一层继续插入
	child.insert(pattern, parts, height+1)
}

// search 查找节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 到达叶子节点 或者 匹配到“*”
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		// 在下一层中继续查找
		result := child.search(parts, height+1)
		if result != nil {
			// 找到节点
			return result
		}
	}

	return nil
}
