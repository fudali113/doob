package tree_router

import (
	"log"
	"testing"

	"github.com/fudali113/doob/core/router"
)

func TestNode_insertChildren(t *testing.T) {

	node := &node{
		class:    normal,
		value:    nil,
		handler:  nil,
		children: make([]*node, 0),
	}

	testUrl := "/oooo/lllll/ddddd"
	node.insertChild(testUrl, &router.SimpleRestHandler{})
	res, err := node.getRT(testUrl, nil)
	log.Print(res)
	if err != nil || res == nil {
		t.Error("Node_insertChildren have bug")
	}

	testUrl1 := "/oooo/bbbb/**"
	_testUrl1 := "/oooo/bbbb/o"
	node.insertChild(testUrl1, &router.SimpleRestHandler{})
	_, err1 := node.getRT(_testUrl1, nil)
	if err1 != nil {
		t.Error("Node_insertChildren have bug 1")
	}

	testUrl2 := "hhhh/{mmm:\\d{3}}/ddddd"
	_testUrl2 := "hhhh/124/ddddd"
	paramMap := map[string]string{}
	node.insertChild(testUrl2, &router.SimpleRestHandler{})
	_, err2 := node.getRT(_testUrl2, paramMap)
	if err2 != nil {
		t.Error("Node_insertChildren have bug 3")
	}
	if paramMap["mmm"] != "124" {
		t.Error("Node_insertChildren have bug 3")
	}

}

func TestNode_Sort(t *testing.T) {

	nodeRoot := &node{
		class:    normal,
		value:    nil,
		handler:  nil,
		children: make([]*node, 0),
	}

	nodeSlice1 := []*node{
		&node{class: 12},
		&node{class: 15},
		&node{class: 11},
		&node{class: 16},
	}

	nodeSlice2 := []*node{
		&node{class: 16},
		&node{class: 11},
	}
	nodeSlice := []*node{
		&node{class: 0},
		&node{class: 1},
		&node{class: 2, children: nodeSlice2},
		&node{class: 3},
		&node{class: 2},
		&node{class: 5, children: nodeSlice1},
		&node{class: 1},
		&node{class: 6},
	}

	nodeRoot.children = nodes(nodeSlice)
	nodeRoot.Sort()
}
