package main

import (
	"net/http"

	"github.com/fudali113/doob"
)

/**
 * 开始http服务
 */
func main() {
	doob.AddHandlerFunc("/doob/{who}/{do}", test1, doob.GET, doob.POST, doob.PUT, doob.DELETE)
	doob.Get("/doob/{haha}", func(w http.ResponseWriter, r *http.Request) {
		haha := r.Form.Get("haha")
		w.Write([]byte(haha))
	})
	doob.Start(8888)
}

func test1(w http.ResponseWriter, r *http.Request) {
	who := r.Form.Get("who")
	do := r.Form.Get("do")
	w.Write([]byte(who + ":" + do))
}
