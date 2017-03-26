package router

// import (
// 	"log"
// 	"testing"
// )

// var (
// 	testNode = &Node{
// 		class:    normal,
// 		value:    nil,
// 		handler:  nil,
// 		children: new(childrens),
// 	}
// )

// func TestNode_insertChildren(t *testing.T) {

// 	testURL := ""
// 	testNode.InsertChild(testURL, &SimpleRestHandler{})
// 	res, err := testNode.GetRT(testURL, nil)
// 	if err != nil || res == nil {
// 		t.Error("Node_insertChildren have bug")
// 	}

// 	testURL1 := "/oooo/*/*"
// 	_testURL1 := "/oooo/bbbb/o/oooo"
// 	testNode.InsertChild(testURL1, &SimpleRestHandler{})
// 	paramMap := map[string]string{}
// 	_, err1 := testNode.GetRT(_testURL1, paramMap)
// 	log.Println("dddddddddddd", paramMap)
// 	if err1 != nil {
// 		t.Error("Node_insertChildren have bug 1")
// 	}

// 	testURL2 := "hhhh/{mmm:\\d{3}}/ddddd"
// 	_testURL2 := "hhhh/124/ddddd"
// 	paramMap = map[string]string{}
// 	testNode.InsertChild(testURL2, &SimpleRestHandler{})
// 	_, err2 := testNode.GetRT(_testURL2, paramMap)
// 	if err2 != nil {
// 		t.Error("Node_insertChildren have bug 3")
// 	}
// 	log.Println("dddddddddddd", paramMap)
// 	if paramMap["mmm"] != "124" {
// 		t.Error("Node_insertChildren have bug 3 + ")
// 	}
// 	log.Print(testNode)
// }

// func TestNode_getNode(t *testing.T) {
// 	getNode := testNode.GetNode("oooo")
// 	log.Print("-------", getNode)
// }
