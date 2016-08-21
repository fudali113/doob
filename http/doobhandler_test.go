package http

import (
	"testing"
	"net/http"
	"fmt"
)

func TestAddHandlerFunc(t *testing.T) {
	AddHandlerFunc("/oo/aa/bb","get", func(http.ResponseWriter , *http.Request) {
		fmt.Println("oo/aa/bb")
	})
	AddHandlerFunc("/oo/{aa}/cc","get", func(http.ResponseWriter , *http.Request) {
		fmt.Println("oo/{aa}/cc")
	})
	AddHandlerFunc("/oo/aa/bb/*","get", func(http.ResponseWriter , *http.Request) {
		fmt.Println("oo/aa/bb/*")
	})
	t.Error(handlerMap)
}
