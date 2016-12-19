package tree_router

import (
	"sort"

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
	value    string
	handler  router.RestHandler
	children nodes
}

// 插入一个子node到一个node中
// 递归插入
// 知道url到最后
func (this *node) insertChild(url string, rt reserveType) error {
	prefix, other := splitUrl(url)
	for _, node := range this.children {
		if node.value == prefix && node.class == normal {
			node.insertChild(other, rt)
		}
	}
	newNode, isOver := creatNode(url, rt)
	if !isOver {
		newNode.insertChild(other, rt)
	}
	return nil
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
