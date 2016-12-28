//
// 在用户和实际逻辑之间做一个中转
// 封装相关方法帮用户更好的使用
//
package doob

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fudali113/doob/router"
	"github.com/fudali113/doob/utils"
)

// http method
const (
	GET     HttpMethod = "get"
	POST    HttpMethod = "post"
	PUT     HttpMethod = "put"
	DELETE  HttpMethod = "delete"
	OPTIONS HttpMethod = "options"
	HEAD    HttpMethod = "head"

	url_split_symbol = "&&"
)

var (
	staticFileCache = map[string][]byte{}

	filters = make([]Filter, 0)
	root    = router.GetRoot()

	_doob = &doob{
		filters: filters,
		root:    root,
	}
)

// start doob server
func Start(port int) {
	log.Printf("server is starting , listen port is %d", port)
	err := Listen(port)
	if err != nil {
		log.Printf("start is fail => %s", err.Error())
	}
}

func DefaultRouter() Router {
	return Router{node: root}
}

func GetRouter(prefix string) Router {
	node := root.GetNode(prefix)
	return Router{node: node}
}

func Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), _doob)
}

func AddFilter(fs ...Filter) {
	filters = append(filters, fs...)
}

func AddStaticPrefix(prefixs ...string) {
	for _, prefixUrl := range prefixs {
		prefixUrl = prefixUrl + "/**"
		DefaultRouter().AddHandlerFunc(prefixUrl, staticPrefixFileHandlerFunc, GET)
	}
}

// static file handler func
func staticPrefixFileHandlerFunc(w http.ResponseWriter, r *http.Request) {
	path := func(url string) string {
		return strings.TrimPrefix(url, "/")
	}(r.URL.Path)

	if utils.IsDirectory(path) {
		path = path + "index.html"
	}

	fileBytes, ok := staticFileCache[path]
	if ok {
		w.Write(fileBytes)
		return
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	staticFileCache[path] = fileBytes
	w.WriteHeader(200)
	w.Write(fileBytes)
}
