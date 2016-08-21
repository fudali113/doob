package http

import (
	"testing"
	"net/http"
	"fmt"
)

func TestAddHandlerFunc(t *testing.T) {
	AddHandlerFunc("/oo/aa/bb","get,post", func(http.ResponseWriter , *http.Request) {
		fmt.Println("oo/aa/bb")
	})
	AddHandlerFunc("/oo/{aa}/cc","PUT", func(http.ResponseWriter , *http.Request) {
		fmt.Println("oo/{aa}/cc")
	})
	AddHandlerFunc("/oo/aa/bb/*","get", func(http.ResponseWriter , *http.Request) {
		fmt.Println("oo/aa/bb/*")
	})
	fmt.Println(handlerMap.simple["/oo/aa/bb"].methodStr)
	oo := handlerMap.rest.urls
	for k,v:=range oo{
		fmt.Println(k,v)
	}
	handler1 := handlerMap.simple["/oo/aa/bb"]
	if !handler1.matchMethod("Get") {
		t.Errorf("match mothod bug")
	}
	handler2,paras := handlerMap.rest.getHandler("/sssss/eooeoeo/cc")
	if handler2 == nil {
		t.Errorf("match url bug")
	}else{
		if !handler2.matchMethod("put") {
			t.Errorf("match mothod bug")
		}
		fmt.Println(paras)
	}

}