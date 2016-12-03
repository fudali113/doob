package doob

/*
现实思路:
    分为三种类型url
	1.普通的:如/fff/ddd/lll
	2.有映射值的:如/user/{who}/info
	3.尾部全部匹配的:如/user/*('*'只可以用于尾部)
    先进行分组,将普通的于要进行取值的分开
    普通的对于一个map,直接使用map[string]获取handlerFunc
    要取值得将url利用"/"切分成数组,匹配是再将实际url切分成数组进行比对并取值

*/

import (
	"log"
	"net/http"

	"github.com/fudali113/doob/core"
)

const (
	GET     HttpMethod = "get"
	POST    HttpMethod = "post"
	PUT     HttpMethod = "put"
	DELETE  HttpMethod = "delete"
	OPTIONS HttpMethod = "options"
	HEAD    HttpMethod = "head"
)

/**
 * 启动server
 */
func Start(port int) {
	log.Printf("server is starting , listen port is %d", port)
	err := core.Listen(port)
	if err != nil {
		log.Printf("start is fail => %s", err.Error())
	}
}

/**
 * 注册一个handler
 */
func AddHandlerFunc(url string, handler http.HandlerFunc, tms ...core.HttpMethod) {
	core.AddHandlerFunc(url, handler, tms...)
}

func Get(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, handler, GET)
}
func Post(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, handler, POST)
}
func Put(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, handler, PUT)
}
func Delete(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, handler, DELETE)
}

/**
 * 添加一个过滤器
 */
func AddFilter(f core.Filter) {
	core.AddFilter(f)
}

func AddFilters(fs ...core.Filter) {
	for i := 0; i < len(fs); i++ {
		core.AddFilter(fs[i])
	}
}
