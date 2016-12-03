# doob

doob is a rest and a simple router handler

init invoke AddHandlerFunc(url,methodStr,func)
such as
```
doob.AddHandlerFunc("/test/{who}/{do}","get,post", func(w http.ResponseWriter,r *http.Request) {
    who := r.Form.Get("who")
    do := r.Form.Get("do")
    w.Write([]byte(who+":"+do))
})

doob.Get("/test/{who}",func(w http.ResponseWriter,r *http.Request) {
    who := r.Form.Get("who")
    w.Write([]byte(who))
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
