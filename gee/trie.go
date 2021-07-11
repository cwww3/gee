package gee

import "strings"

type node struct {
	pattern   string // 待匹配路由(完整)
	part     string  // 部分
	children []*node
	isWild   bool // 是否精确匹配
}

// 找第一个匹配成的
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// 找到所有的可能
func (n *node) matchChildren(part string) []*node {
	var nodes = make([]*node, 0, len(n.children))
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归插入路径 路由匹配
func (n *node) insert(pattern string, parts []string, height int) {
	// 全部插入完成 会在最有一个节点存储完整pattern
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// 获取当前部分
	part := parts[height]
	child := n.matchChild(part)
	if child == nil { // 说明当前部分还不存在构建一个
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 递归搜索路径
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 插入时在最后一个节点存储了完整路径
		if n.part == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil { // 找到了
			return result
		}
	}
	return nil
}
