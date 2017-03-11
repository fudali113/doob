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

	"github.com/fudali113/doob/http/router"
	"github.com/fudali113/doob/middleware"
	"github.com/fudali113/doob/utils"

	_ "github.com/fudali113/doob/middleware/basicauth"

	. "github.com/fudali113/doob/http/const"
)

var (
	staticFileCache = map[string][]byte{}

	beforeFilters = make([]middleware.BeforeFilter, 0)
	laterFilters  = make([]middleware.LaterFilter, 0)
	root          = router.GetRoot()

	doob = &Doob{
		bFilters: middleware.BFilters,
		lFilters: middleware.LFilters,
		Root:     root,
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
	return http.ListenAndServe(fmt.Sprintf(":%d", port), doob)
}

func AddBFilter(fs ...middleware.BeforeFilter) {
	beforeFilters = append(beforeFilters, fs...)
}

func AddLFilter(fs ...middleware.LaterFilter) {
	laterFilters = append(laterFilters, fs...)
}

func AddMiddlerware(fs ...middleware.Middleware) {
	middleware.AddMiddleware(fs...)
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
		w.WriteHeader(NOT_FOUND)
		return
	}

	staticFileCache[path] = fileBytes
	w.WriteHeader(OK)
	w.Write(fileBytes)
}

func init() {
	AddLFilter(middleware.HeadHTTPMethodDealer(func(w http.ResponseWriter, req *http.Request) {
		methodStr := strings.ToLower(req.Method)
		if methodStr == string(HEAD) {
			w.Write([]byte{})
		}
	}))
}
