package router

import (
	"log"
	"testing"
)

var (
	testNode = &Node{
		class:    normal,
		value:    nil,
		handler:  nil,
		children: make([]*Node, 0),
	}
)

func TestNode_insertChildren(t *testing.T) {

	testUrl := ""
	testNode.InsertChild(testUrl, &SimpleRestHandler{})
	res, err := testNode.GetRT(testUrl, nil)
	if err != nil || res == nil {
		t.Error("Node_insertChildren have bug")
	}

	testUrl1 := "/oooo/bbbb/**"
	_testUrl1 := "/oooo/bbbb/o"
	testNode.InsertChild(testUrl1, &SimpleRestHandler{})
	_, err1 := testNode.GetRT(_testUrl1, nil)
	if err1 != nil {
		t.Error("Node_insertChildren have bug 1")
	}

	testUrl2 := "hhhh/{mmm:\\d{3}}/ddddd"
	_testUrl2 := "hhhh/124/ddddd"
	paramMap := map[string]string{}
	testNode.InsertChild(testUrl2, &SimpleRestHandler{})
	_, err2 := testNode.GetRT(_testUrl2, paramMap)
	if err2 != nil {
		t.Error("Node_insertChildren have bug 3")
	}
	if paramMap["mmm"] != "124" {
		t.Error("Node_insertChildren have bug 3")
	}
	log.Print(testNode)
}

func TestNode_Sort(t *testing.T) {

	nodeRoot := &Node{
		class:    normal,
		value:    nil,
		handler:  nil,
		children: make([]*Node, 0),
	}

	nodeSlice1 := []*Node{
		&Node{class: 12},
		&Node{class: 15},
		&Node{class: 11},
		&Node{class: 16},
	}

	nodeSlice2 := []*Node{
		&Node{class: 16},
		&Node{class: 11},
	}
	nodeSlice := []*Node{
		&Node{class: 0},
		&Node{class: 1},
		&Node{class: 2, children: nodeSlice2},
		&Node{class: 3},
		&Node{class: 2},
		&Node{class: 5, children: nodeSlice1},
		&Node{class: 1},
		&Node{class: 6},
	}

	nodeRoot.children = nodes(nodeSlice)
	nodeRoot.Sort()
}

func TestNode_getNode(t *testing.T) {
	getNode := testNode.GetNode("oooo")
	log.Print("-------", getNode)
}
