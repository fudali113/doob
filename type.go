package doob

import (
	"log"
	"net/http"

	"github.com/fudali113/doob/router"
	"github.com/fudali113/doob/utils"
)

// Filter接口
type Filter interface {
	// Filter 的实际操作
	// 返回 bool 值决定是否通过此 filter
	doFilter(res http.ResponseWriter, req *http.Request) bool
}

// 封装node，对外提供简单方法
type Router struct {
	node *router.Node
}

func (r Router) AddHandlerFunc(allUrl string, handler interface{}, methods ...HttpMethod) {
	urls := utils.Split(allUrl, urlSplitSymbol)
	for _, url := range urls {
		for _, method := range methods {
			methodStr := string(method)
			if checkMethodStr(methodStr) {
				r.node.InsertChild(url, router.GetSimpleRestHandler(methodStr, handler))
			} else {
				log.Printf("%s method is unsupport", methodStr)
			}
		}
	}
}

func (r Router) Get(allUrl string, handler interface{}) {
	r.AddHandlerFunc(allUrl, handler, GET)
}

func (r Router) Post(allUrl string, handler interface{}) {
	r.AddHandlerFunc(allUrl, handler, POST)
}

func (r Router) Put(allUrl string, handler interface{}) {
	r.AddHandlerFunc(allUrl, handler, PUT)
}

func (r Router) Delete(allUrl string, handler interface{}) {
	r.AddHandlerFunc(allUrl, handler, DELETE)
}

func (r Router) Options(allUrl string, handler interface{}) {
	r.AddHandlerFunc(allUrl, handler, OPTIONS)
}

func (r Router) Head(allUrl string, handler interface{}) {
	r.AddHandlerFunc(allUrl, handler, HEAD)
}
