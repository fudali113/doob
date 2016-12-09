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
	doob.Get("/doob/test/{name}/{value}", test)
	doob.Get("/doob/{haha:[0-9]{3,4}}/ooo/kkkk", func(w http.ResponseWriter, r *http.Request) {
		haha := r.Form.Get("haha")
		w.Write([]byte(haha))
	})
	doob.Start(8888)
}

func test(name, value string) map[string]string {
	return map[string]string{
		value: name,
	}
}

func test1(w http.ResponseWriter, r *http.Request) {
	who := r.Form.Get("who")
	do := r.Form.Get("do")
	w.Write([]byte(who + ":" + do))
}
