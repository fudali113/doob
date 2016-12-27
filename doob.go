//
// 在用户和实际逻辑之间做一个中转
// 封装相关方法帮用户更好的使用
//
package doob

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fudali113/doob/core"
	"github.com/fudali113/doob/log"
	"github.com/fudali113/doob/utils"
)

// http method
const (
	GET     core.HttpMethod = "get"
	POST    core.HttpMethod = "post"
	PUT     core.HttpMethod = "put"
	DELETE  core.HttpMethod = "delete"
	OPTIONS core.HttpMethod = "options"
	HEAD    core.HttpMethod = "head"
)

var (
	logger          = log.GetLog("doob")
	staticFileCache = map[string][]byte{}
)

// start doob server
func Start(port int) {
	logger.Info("server is starting , listen port is %d", port)
	err := core.Listen(port)
	if err != nil {
		logger.Error("start is fail => %s", err.Error())
	}
}

// register a handler
func AddHandlerFunc(url string, handler interface{}, tms ...core.HttpMethod) {
	core.AddHandlerFunc(url, handler, tms...)
}

func Get(url string, handler interface{}) {
	core.AddHandlerFunc(url, handler, GET)
}
func Post(url string, handler http.HandlerFunc) {
	core.AddHandlerFunc(url, handler, POST)
}
func Put(url string, handler http.HandlerFunc) {
	core.AddHandlerFunc(url, handler, PUT)
}
func Delete(url string, handler http.HandlerFunc) {
	core.AddHandlerFunc(url, handler, DELETE)
}

// set return type default value
func SetReturnDealDefaultType(t string) {
	core.SetReturnDealDefaultType(t)
}

// add a filter
func AddFilter(fs ...core.Filter) {
	core.AddFilter(fs...)
}

func AddStaicPrefix(prefixs ...string) {
	for _, prefixUrl := range prefixs {
		prefixUrl = prefixUrl + "/**"
		core.AddHandlerFunc(prefixUrl, staticPrefixFileHandlerFunc, GET)
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
