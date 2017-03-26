package router

import (
	"encoding/json"
	"fmt"
	"log"
)

// node的类别
const (
	normal = iota
	pathReg
	pathVar
	matchAll
)

var (
	root = &Node{Class: normal}
)

// GetRoot get a root node
// this node start with '/'
func GetRoot() *Node {
	return root
}

// ReserveType 转化储存的实体类型
// 一面后面修改麻烦
type ReserveType RestHandler

// conflictRTDealFunc 注册时冲突处理函数type
type conflictRTDealFunc func(old, new ReserveType) ReserveType

var (
	crtdf conflictRTDealFunc = func(old, new ReserveType) ReserveType {
		old.Joint(new)
		return old
	}
)

// Node 装载url每一个由`/`隔开的分段的实体
// 递归结构
type Node struct {
	Class    int
	Value    nodeV
	handler  RestHandler
	Children *childrens
}

// InsertChild 插入一个子node到一个node中
// 递归插入
// 直到url到最后
func (n *Node) InsertChild(URL string, rt ReserveType) {
	n.insert(URL, rt)
}

func (n *Node) insert(URL string, rt ReserveType) {
	prefix, other := splitURL(URL)
	if prefix == "" {
		if n.handler == nil {
			n.handler = rt
		}
		n.handler = crtdf(n.handler, rt)
		return
	}
	if n.Children == nil {
		n.Children = new(childrens)
	}
	node := n.Children.getNode(prefix, nil)
	if node == nil {
		node = creatNode(prefix)
	}
	node.insert(other, rt)
	n.Children.insert(node)

}

// creatNode create a new Node
func creatNode(passageURL string) *Node {
	return &Node{
		Class:    getClass(passageURL),
		Value:    createNodeValue(passageURL),
		Children: new(childrens),
	}
}

// GetRT get reserve type
// if reserve type value is nil , return error
func (n *Node) GetRT(url string, paramMap map[string]string) (ReserveType, error) {
	log.Println("root_>", root)
	isMatch := false
	node := n.getNode(url, paramMap, &isMatch)
	if !isMatch {
		return nil, NotMatch{"this url not rt"}
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

	prefix, other := splitURL(url)
	if prefix == "" {
		setNode(n)
		return
	}
	childrens := n.Children
	defer func() {
		if childrens.SuffixMatch != nil && !*isMatch {
			addValueToPathParam(paramMap, suffixMatchSymbol, url)
			setNode(childrens.SuffixMatch)
		}
	}()
	node = childrens.getNode(prefix, paramMap)
	if node != nil {
		if other == "" {
			setNode(node)
		} else {
			node = node.getNode(other, paramMap, isMatch)
		}
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
	return n.getNode(url, nil, &_true)
}

// String 打印内容
func (n *Node) String() string {
	jsonStr, _ := json.Marshal(n)
	return string(jsonStr)
}

// childrens 用于封装该node的所有子node
// 不同的类型使用不同的储存方式
// 以提高性能
type childrens struct {
	Normal   map[string]*Node
	Regexp   []*Node
	AllMatch *Node
	// 尾部全匹配以栈的形式随方法存入
	// 当最后没有匹配是，将获取栈中的倒数第一个元素放回
	SuffixMatch *Node
}

// getNode 根据以`/`分段的url获取子Node
// passageURL 一段URL
func (c *childrens) getNode(passageURL string, paramMap map[string]string) (node *Node) {
	ok := false
	var v *Node
	if c.Normal != nil {
		v, ok = c.Normal[passageURL]
	}
	if ok {
		node = v
	} else if c.Regexp != nil {
		for i := 0; i < len(c.Regexp); i++ {
			nowNode := c.Regexp[i]
			value := nowNode.Value
			if match := value.isMatch(passageURL); match {
				node = nowNode
				value.paramValue(passageURL, paramMap)
				break
			}
		}
	} else if c.AllMatch != nil {
		node = c.AllMatch
		c.AllMatch.Value.paramValue(passageURL, paramMap)
	}
	return
}

func (c *childrens) insert(node *Node) {
	switch node.Class {
	case normal:
		if c.Normal == nil {
			c.Normal = make(map[string]*Node)
		}
		c.Normal[node.Value.getOrigin()] = node
	case pathReg:
		if c.Regexp == nil {
			c.Regexp = make([]*Node, 0, 3)
		}
		c.Regexp = append(c.Regexp, node)
	case pathVar:
		c.AllMatch = node
	case matchAll:
		c.SuffixMatch = node
	}
}

func (c *childrens) String() string {
	return fmt.Sprintf("{ normal:%v,regexp:%v,allMatch:%v,suffixMatch:%v }", c.Normal, c.Regexp, c.AllMatch, c.SuffixMatch)
}
