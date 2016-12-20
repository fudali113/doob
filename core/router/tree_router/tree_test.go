package tree_router

import (
	"testing"
	"github.com/fudali113/doob/core/router"
)

func TestNode_insertChildren(t *testing.T) {
	testUrl := "/oooo/lllll/ddddd"
	node := &node{
		class:normal,
		value:"",
		handler:nil,
		children:make([]*node,0),
	}
	node.insertChild(testUrl,&router.SimpleRestHandler{})
	_ ,err := node.getRT(testUrl)
	if err != nil {
		t.Error("Node_insertChildren have bug")
	}

}
