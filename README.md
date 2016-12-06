# doob

doob is a rest and a simple router handler

init invoke AddHandlerFunc(url,methodStr,func)
such as
```
func test1(w http.ResponseWriter, r *http.Request) {
	who := r.Form.Get("who")
	do := r.Form.Get("do")
	w.Write([]byte(who + ":" + do))
}

doob.AddHandlerFunc("/doob/{who}/{do}", test1, doob.GET, doob.POST, doob.PUT, doob.DELETE)
doob.Get("/doob/{haha}", func(w http.ResponseWriter, r *http.Request) {
  haha := r.Form.Get("haha")
  w.Write([]byte(haha))
})

//use regexp in your url template
doob.Get("/ooo/{name:\\w+}",func(w http.ResponseWriter, r *http.Request) {
  haha := r.Form.Get("name")
  w.Write([]byte(haha))
})
```
use r.From.Get() receive your urlpara value

next use
```
doob.Start(8888)
```
run you application

clone this project , run demo
```
cd /demo
go run demo.go
```
