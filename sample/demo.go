package main

import (
	"net/http"

	"github.com/fudali113/doob"
)

// 开始http服务
func main() {
	doob.AddStaicPrefix("/static")
	doob.AddHandlerFunc("/doob/origin/{who}/{do}", origin, doob.GET, doob.POST, doob.PUT, doob.DELETE)
	doob.Get("/doob/di/{name}/{value}", di)
	doob.Get("/doob/ctx/{haha:[0-9]{3,4}}", ctx)
	doob.Get("/test", returnHtml)
	doob.Start(8888)
}

// 根据url参数自动注入参数
// 返回值为string时为返回静态文件
// 返回不为string时默认将解析该对象，并返回给请求用户
func di(name, value string) map[string]string {
	return map[string]string{
		value: name,
	}
}

// 兼容原始http方法类
func origin(w http.ResponseWriter, r *http.Request) {
	who := r.Form.Get("who")
	do := r.Form.Get("do")
	w.Write([]byte(who + ":" + do))
}

// 根据doob 里的context 进行获取参数或者返回
func ctx(ctx *doob.Context) interface{} {
	return map[string]int{
		"haha": ctx.ParamInt("haha"),
		"test": ctx.ParamInt("test"),
	}
}

// 返回处理 template 文件 path 和数据进行处理并返回生成的html
func returnHtml() (string, interface{}) {
	return "tpl:static/test", map[string]string{"Name": "sdddddddddddddddddddddddddddddd"}
}
