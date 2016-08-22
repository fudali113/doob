# golib

DoobHandler is a rest handler

init invoke AddHandlerFunc(url,methodStr,func)
such as
`AddHandlerFunc("/test/{who}/{do}","get,post", func(w http.ResponseWriter,r *http.Request) {
    who := r.Form.Get("who")
    do := r.Form.Get("do")
    w.Write([]byte(who+":"+do))
})`

use r.From.Get() receive your urlpara value