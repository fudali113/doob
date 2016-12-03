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
func AddHandlerFunc(url string, methodStr string, handler http.HandlerFunc) {
	core.AddHandlerFunc(url, methodStr, handler)
}

func Get(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, "get", handler)
}
func Post(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, "post", handler)
}
func Put(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, "put", handler)
}
func Delete(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, "delete", handler)
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
