package tree_router

import (
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
	_, err := node.getRT(testUrl, nil)
	if err != nil {
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
