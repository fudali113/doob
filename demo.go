package main

import (
	myhandler "github.com/fudali113/golib/http"
	"net/http"
)
func main()  {
	myhandler.AddHandlerFunc("/test/{who}/{do}","get,post", func(w http.ResponseWriter,r *http.Request) {
		who := r.Form.Get("who")
		do := r.Form.Get("do")
		w.Write([]byte(who+":"+do))
	})
	http.ListenAndServe(":3333",myhandler.DoobHandler{})
}
