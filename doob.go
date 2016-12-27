//
// 在用户和实际逻辑之间做一个中转
// 封装相关方法帮用户更好的使用
//
package doob

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
)

var (
	staticFileCache = map[string][]byte{}
)

// start doob server
func Start(port int) {
	log.Printf("server is starting , listen port is %d", port)
	err := Listen(port)
	if err != nil {
		log.Printf("start is fail => %s", err.Error())
	}
}

func Get(url string, handler interface{}) {
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
func Options(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, handler, OPTIONS)
}
func Head(url string, handler http.HandlerFunc) {
	AddHandlerFunc(url, handler, HEAD)
}

func AddStaicPrefix(prefixs ...string) {
	for _, prefixUrl := range prefixs {
		prefixUrl = prefixUrl + "/**"
		AddHandlerFunc(prefixUrl, staticPrefixFileHandlerFunc, GET)
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
