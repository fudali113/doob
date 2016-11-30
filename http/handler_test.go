package http

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAddHandlerFunc(t *testing.T) {
	AddHandlerFunc("/oo/aa/bb", "get,post", func(http.ResponseWriter, *http.Request) {
		fmt.Println("oo/aa/bb")
	})
	AddHandlerFunc("/oo/{aa}/cc/{bb}", "get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Form.Get("aa"), r.Form.Get("bb"))
	})
	AddHandlerFunc("/oo/aa/bb/*", "get", func(http.ResponseWriter, *http.Request) {
		fmt.Println("oo/aa/bb/*")
	})
	fmt.Println(handlerMap.simple["/oo/aa/bb"].methodStr)
	oo := handlerMap.rest.urls
	for k, v := range oo {
		fmt.Println(k, v)
	}
	handler1 := handlerMap.simple["/oo/aa/bb"]
	if !handler1.matchMethod("Get") {
		t.Errorf("match mothod bug")
	}
	handler2, paras := handlerMap.rest.getHandler("/oo/eooeoeo/cc/kkkkkkkk")
	if handler2 == nil {
		t.Errorf("match url bug")
	} else {
		if !handler2.matchMethod("GET") {
			t.Errorf("match mothod bug")
		}
		fmt.Println(paras)
	}
	fmt.Println("--->", handlerMap.lastAllMatch)
	handler3, paras := handlerMap.rest.getHandler("/oo/aa/bb/123456")
	if handler3 == nil {
		t.Errorf("match url bug")
	} else {
		if !handler3.matchMethod("get") {
			t.Errorf("match mothod bug")
		}
		fmt.Println(paras)
	}
	req, err := http.NewRequest("GET", "http://localhost:8080/oo/aa/bb/123456/4232324342423", nil)
	req.Form = map[string][]string{}
	if err != nil {
		return
	}
	handler, err1 := handlerMap.getHandler(req)
	var handlerFunc func(http.ResponseWriter, *http.Request) = handler
	if err1 != nil {
		t.Error(err1)
	} else {
		fmt.Println("all")
		handlerFunc(nil, req)
	}

}
