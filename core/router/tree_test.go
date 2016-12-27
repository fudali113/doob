package router

import (
	"log"
	"testing"
)

var (
	testNode = &node{
		class:    normal,
		value:    nil,
		handler:  nil,
		children: make([]*node, 0),
	}
)

func TestNode_insertChildren(t *testing.T) {

	testUrl := ""
	testNode.insertChild(testUrl, &SimpleRestHandler{})
	res, err := testNode.getRT(testUrl, nil)
	if err != nil || res == nil {
		t.Error("Node_insertChildren have bug")
	}

	testUrl1 := "/oooo/bbbb/**"
	_testUrl1 := "/oooo/bbbb/o"
	testNode.insertChild(testUrl1, &SimpleRestHandler{})
	_, err1 := testNode.getRT(_testUrl1, nil)
	if err1 != nil {
		t.Error("Node_insertChildren have bug 1")
	}

	testUrl2 := "hhhh/{mmm:\\d{3}}/ddddd"
	_testUrl2 := "hhhh/124/ddddd"
	paramMap := map[string]string{}
	testNode.insertChild(testUrl2, &SimpleRestHandler{})
	_, err2 := testNode.getRT(_testUrl2, paramMap)
	if err2 != nil {
		t.Error("Node_insertChildren have bug 3")
	}
	if paramMap["mmm"] != "124" {
		t.Error("Node_insertChildren have bug 3")
	}
	log.Print(testNode)
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

func TestNode_getNode(t *testing.T) {
	getNode := testNode.getNode("oooo")
	log.Print("-------", getNode)
}
