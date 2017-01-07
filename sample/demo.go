package main

import (
	"net/http"

	"github.com/fudali113/doob"

	. "github.com/fudali113/doob/http/const"

	doobhttp "github.com/fudali113/doob/http"
)

// 开始http服务
func main() {
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
func ctx(ctx *doobhttp.Context) interface{} {
	return map[string]interface{}{
		"haha": ctx.ParamInt("haha"),
		"test": ctx.Seesion().Get("test"),
	}
}

func testForward(ctx *doobhttp.Context) {
	ctx.Forward("", "http://www.23mofang.com")
}

func testRedirect(ctx *doobhttp.Context) {
	ctx.Redirect("test")
}

// 返回处理 template 文件 path 和数据进行处理并返回生成的html
func returnHtml(ctx *doobhttp.Context) (string, interface{}) {
	session := ctx.Seesion()
	session.Set("test", "ooooooooo")
	return "tpl:static/test", map[string]string{"Name": "sdddddddddddddddddddddddddddddd"}
}

// init router
func init() {
	doob.AddStaticPrefix("/static")
	router := doob.DefaultRouter()
	router.AddHandlerFunc("/doob/origin/{who}/{do}", origin, GET, POST, PUT, DELETE)
	router.Get("forward", testForward)
	router.Get("redirect1", testRedirect)
	router.Get("/test", returnHtml)
	doobRouter := doob.GetRouter("doob")
	doobRouter.Get("/di/{name}/{value}", di)
	doobRouter.Get("/ctx/{haha:[0-9]{3,4}}", ctx)
}
