package router

import "fmt"

const (
	normal = iota
	pathReg
	pathVar
	matchAll
)

// GetRoot get a root node
// this node start with '/'
func GetRoot() *Node {
	return &Node{class: normal}
}

// ReserveType 转化储存的实体类型
// 一面后面修改麻烦
type ReserveType RestHandler

// Node 装载url每一个由`/`隔开的分段的实体
// 递归结构
type Node struct {
	class    int
	value    nodeV
	handler  RestHandler
	children *childrens
}

// InsertChild 插入一个子node到一个node中
// 递归插入
// 直到url到最后
func (n *Node) InsertChild(URL string, rt ReserveType) {
	n.insert(URL, rt)
}

func (n *Node) insert(URL string, rt ReserveType) {
	prefix, other := splitUrl(URL)
	if prefix == "" {
		if n.handler == nil {
			n.handler = rt
		}
		n.handler.Joint(rt)
		return
	}
	if n.children == nil {
		n.children = new(childrens)
	}
	node := n.children.getNode(prefix)
	if node == nil {
		node = creatNode(prefix)
	}
	node.insert(other, rt)
	n.children.insert(node)

}

// GetRT get reserve type
// if reserve type value is nil , return error
func (n *Node) GetRT(url string, paramMap map[string]string) (ReserveType, error) {
	isMatch := false
	node := n.getNode(url, paramMap, &isMatch)
	if !isMatch {
		return nil, fmt.Errorf("not match")
	}
	return node.handler, nil
}

// getNode 递归获取 匹配的node
// 用户获取路由Node
// 与下方的GetNode不同，GetNode用于获取注册Node
func (n *Node) getNode(url string, paramMap map[string]string, isMatch *bool) (node *Node) {

	// 设置node的值并返回
	var setNode = func(n *Node) {
		*isMatch = true
		node = n
	}

	prefix, other := splitUrl(url)
	if prefix == "" {
		setNode(n)
		return
	}
	childrens := n.children
	defer func() {
		if childrens.suffixMatch != nil && !*isMatch {
			setNode(childrens.suffixMatch)
		}
	}()
	node = childrens.getNode(prefix)
	if other == "" {
		if node != nil {
			setNode(node)
		}
	} else {
		node = node.getNode(other, paramMap, isMatch)
	}
	return
}

// GetNode 根据url 获取一个注册Node
func (n *Node) GetNode(url string) *Node {
	if url == "" {
		return n
	}
	_, err := n.GetRT(url, nil)
	if err != nil {
		n.InsertChild(url, nil)
	}
	_true := true
	return n.getNode(url, map[string]string{}, &_true)
}

func (n *Node) String() string {
	return fmt.Sprintf("{ class:%d,value:%v,handler:%v,children:%v }", n.class, n.value, n.handler, n.children)
}

// creatNode create a new Node
func creatNode(passageURL string) *Node {
	return &Node{
		class:    getClass(passageURL),
		value:    createNodeValue(passageURL),
		children: new(childrens),
	}
}

func getRtAndErr(rt ReserveType) (ReserveType, error) {
	if rt == nil {
		return nil, NotMatch{"this url not rt"}
	}
	return rt, nil
}
