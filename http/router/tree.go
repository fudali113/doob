package router

import (
	"fmt"
	"sort"
	"strings"
)

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
	children nodes
}

// InsertChild 插入一个子node到一个node中
// 递归插入
// 知道url到最后
func (n *Node) InsertChild(url string, rt ReserveType) error {
	prefix, other := splitUrl(url)
	for _, node := range n.children {
		if node.value.getOrigin() == prefix {
			if other == "" {
				if node.handler == nil {
					node.handler = rt
				} else {
					node.handler.Joint(rt)
				}
			} else {
				node.InsertChild(other, rt)
			}
			return nil
		}
	}
	newNode, isOver := creatNode(url, rt)
	n.children = append(n.children, newNode)
	if !isOver {
		newNode.InsertChild(other, rt)
	}
	return nil
}

// GetRT get reserve type
// if reserve type value is nil , return error
func (n *Node) GetRT(url string, paramMap map[string]string) (ReserveType, error) {
	prefix, other := splitUrl(url)
	for _, node := range n.children {
		nodeValue := node.value
		if match, over := nodeValue.isMatch(prefix); match {
			hasParam, paramMapPart := nodeValue.paramValue(prefix, url)
			if hasParam && paramMap != nil {
				for k, v := range paramMapPart {
					paramMap[k] = v
				}
			}
			if over || other == "" {
				return getRtAndErr(node.handler)
			}
			return node.GetRT(other, paramMap)
		}
	}
	return getRtAndErr(nil)
}

// GetNode 根据url 获取一个node
func (n *Node) GetNode(url string) *Node {
	if url == "" {
		return n
	}
	_, err := n.GetRT(url, nil)
	if err != nil {
		n.InsertChild(url, nil)
	}
	prefix, other := splitUrl(url)
	for _, node := range n.children {
		if node.value.getOrigin() == prefix {
			return node.GetNode(other)
		}
	}
	return nil
}

// Sort 对子node进行排序
// 将会递归所有子node排序
func (n *Node) Sort() {
	if n.children == nil {
		return
	}
	sort.Sort(n.children)
	for _, node := range n.children {
		node.Sort()
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("{ class:%d,value:%s,handler:%v,children:%v }", n.class, n.value, n.handler, n.children)
}

// create a new Node
func creatNode(url string, rt ReserveType) (newNode *Node, isOver bool) {
	prefix, other := splitUrl(url)
	isOver = false
	newNode = &Node{
		class: getClass(prefix),
		value: createNodeValue(prefix),
	}
	if strings.TrimSpace(other) == "" {
		newNode.handler = rt
		isOver = true
	} else {
		newNode.children = make([]*Node, 0)
	}
	return
}

func getRtAndErr(rt ReserveType) (ReserveType, error) {
	if rt == nil {
		return nil, NotMatch{"this url not rt"}
	}
	return rt, nil
}
