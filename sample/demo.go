package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fudali113/doob"
	"github.com/fudali113/doob/core"
)

type oo func(string, string)

func (o oo) ooo(name string, value string) {
	o(name, value)
}

func ooo(name string, value string) {
	return
}

func oooo(w http.ResponseWriter, req *http.Request) {
	return
}

/**
 * 开始http服务
 */
func main() {
	_, err := os.Open("oooo.html")
	if err != nil {
		log.Print(err)
	}
	doob.AddStaicPrefix("/static")
	doob.AddHandlerFunc("/doob/origin/{who}/{do}", origin, doob.GET, doob.POST, doob.PUT, doob.DELETE)
	doob.Get("/doob/di/{name}/{value}", di)
	doob.Get("/doob/ctx/{haha:[0-9]{3,4}}", ctx)
	doob.Start(8888)
}

/**
 * 根据url参数自动注入参数
 * 返回值为string时为返回静态文件
 * 返回不为string时默认将解析该对象，并返回给请求用户
 */
func di(name, value string) map[string]string {
	return map[string]string{
		value: name,
	}
}

/**
 * 兼容原始http方法类
 */
func origin(w http.ResponseWriter, r *http.Request) {
	who := r.Form.Get("who")
	do := r.Form.Get("do")
	w.Write([]byte(who + ":" + do))
}

/**
 * 根据doob 里的context 进行获取参数或者返回
 */
func ctx(ctx *core.Context) interface{} {
	return map[string]int{
		"haha": ctx.ParamInt("haha"),
		"test": ctx.ParamInt("test"),
	}
}
