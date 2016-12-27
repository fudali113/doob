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

func GetRoot() *Node {
	return &Node{class: normal}
}

// 转化储存的实体类型
// 一面后面修改麻烦
type reserveType RestHandler

// 装载url每一个由`/`隔开的分段的实体
type Node struct {
	class    int
	value    nodeV
	handler  RestHandler
	children nodes
}

// 插入一个子node到一个node中
// 递归插入
// 知道url到最后
func (this *Node) InsertChild(url string, rt reserveType) error {
	prefix, other := splitUrl(url)
	for _, node := range this.children {
		if node.value.getOrigin() == prefix {
			node.InsertChild(other, rt)
			return nil
		}
	}
	newNode, isOver := creatNode(url, rt)
	this.children = append(this.children, newNode)
	if !isOver {
		newNode.InsertChild(other, rt)
	}
	return nil
}

// get reserve type
// if reserve type value is nil , return error
func (this *Node) GetRT(url string, paramMap map[string]string) (reserveType, error) {
	prefix, other := splitUrl(url)
	for _, node := range this.children {
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
	return nil, NotMatch{"this url not rt"}
}

func (this *Node) GetNode(url string) *Node {
	if url == "" {
		return this
	}
	_, err := this.GetRT(url, nil)
	if err != nil {
		this.InsertChild(url, nil)
	}
	prefix, other := splitUrl(url)
	for _, node := range this.children {
		if node.value.getOrigin() == prefix {
			return node.GetNode(other)
		}
	}
	return nil
}

// 对子node进行排序
// 将会递归所有子node排序
func (this *Node) Sort() {
	if this.children == nil {
		return
	}
	sort.Sort(this.children)
	for _, node := range this.children {
		node.Sort()
	}
}

func (this *Node) String() string {
	return fmt.Sprintf("{ class:%d,value:%s,handler:%v,children:%v }", this.class, this.value, this.handler, this.children)
}

// create a new Node
func creatNode(url string, rt reserveType) (newNode *Node, isOver bool) {
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

func getRtAndErr(rt reserveType) (reserveType, error) {
	if rt == nil {
		return nil, NotMatch{"this url not rt"}
	}
	return rt, nil
}
