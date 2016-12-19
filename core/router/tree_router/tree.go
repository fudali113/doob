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

type nodes []*node

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

type node struct {
	class    int
	value    string
	handler  router.RestHandler
	children nodes
}

func (this *node) Add(url string , rest router.RestHandler) error {

}

func (this *node) Sort() {
	if this.children == nil {
		return
	}
	sort.Sort(this.children)
	for _, node := range this.children {
		node.Sort()
	}
}


