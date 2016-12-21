package tree_router

import (
	"log"
	"sort"
	"strings"

	"fmt"

	"github.com/fudali113/doob/core/router"
)

const (
	normal = iota
	pathReg
	pathVar
	matchAll
)

// 实现Sort的接口
type nodes []*node

// 转化储存的实体类型
// 一面后面修改麻烦
type reserveType router.RestHandler

func (this nodes) Len() int {
	return len(this)
}
func (this nodes) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this nodes) Less(i, j int) bool {
	a := this[i]
	b := this[j]
	return a.class > b.class
}

// 装载url每一个由`/`隔开的分段的实体
type node struct {
	class    int
	value    nodeV
	handler  router.RestHandler
	children nodes
}

// 插入一个子node到一个node中
// 递归插入
// 知道url到最后
func (this *node) insertChild(url string, rt reserveType) error {
	prefix, other := splitUrl(url)
	for _, node := range this.children {
		if node.value.getOrigin() == prefix {
			node.insertChild(other, rt)
			return nil
		}
	}
	newNode, isOver := creatNode(url, rt)
	this.children = append(this.children, newNode)
	if !isOver {
		newNode.insertChild(other, rt)
	}
	return nil
}

// create a new node
func creatNode(url string, rt reserveType) (newNode *node, isOver bool) {
	prefix, other := splitUrl(url)
	isOver = false
	newNode = &node{
		class: getClass(prefix),
		value: createNodeValue(prefix),
	}
	if strings.TrimSpace(other) == "" {
		newNode.handler = rt
		isOver = true
	} else {
		newNode.children = make([]*node, 0)
	}
	return
}

func (this *node) getRT(url string) (reserveType, error) {
	prefix, other := splitUrl(url)
	for _, node := range this.children {
		log.Print(node.value.isMatch(prefix))
		if match, over := node.value.isMatch(prefix); match {
			if over || other == "" {
				return getRtAndErr(node.handler)
			}
			return node.getRT(other)
		}
	}
	return nil, NotMatch{"this url not rt"}
}

func getRtAndErr(rt reserveType) (reserveType, error) {
	if rt == nil {
		return nil, NotMatch{"this url not rt"}
	}
	return rt, nil
}

// 对子node进行排序
// 将会递归所有子node排序
func (this *node) Sort() {
	if this.children == nil {
		return
	}
	sort.Sort(this.children)
	for _, node := range this.children {
		node.Sort()
	}
}

func (this *node) String() string {
	return fmt.Sprintf("{ class:%d,value:%s,handler:%v,children:%v }", this.class, this.value, this.handler, this.children)
}
