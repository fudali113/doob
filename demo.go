package main

import (
	"net/http"

	"github.com/fudali113/golib/doob"
)

/**
 * 开始http服务
 */
func main() {
	doob.AddHandlerFunc("/test/{who}/{do}", "get,post", func(w http.ResponseWriter, r *http.Request) {
		who := r.Form.Get("who")
		do := r.Form.Get("do")
		w.Write([]byte(who + ":" + do))
	})
	doob.AddHandlerFunc("/test/*", "get,post", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Form.Get("*")))
	})
	doob.AddHandlerFunc("", "get", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("oo"))
	})
	doob.Start(8888)
}
