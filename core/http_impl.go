package core

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/fudali113/doob/core/router"
	"github.com/fudali113/doob/log"
	"github.com/fudali113/doob/utils"
)

var (
	logger = log.GetLog("simple router")
)

type doob struct {
	router  router.Router
	filters []Filter
}

func (this *doob) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer logger.Info("程序处理共消耗:%d ns", time.Now().Sub(startTime).Nanoseconds())

	for i := range this.filters {
		if this.filters[i].doFilter(w, req) {
			continue
		} else {
			return
		}
	}

	url := req.URL.Path
	for _, urlPrefix := range urlPrefixs {
		if strings.HasPrefix(url, urlPrefix) {
			serveFile(w, getPath(url))
			return
		}
	}

	matchResult := this.router.Get(url)
	invoke(matchResult, w, req)

}

func (this *doob) addFilter(fs ...Filter) {
	this.filters = append(this.filters, fs...)
}

func (this *doob) addRestHandler(url string, restHandler router.RestHandler) {
	this.router.Add(url, restHandler)
}

func getPath(url string) string {
	if strings.HasPrefix(url, "/") {
		return "." + url
	}
	return "./" + url
}

func serveFile(w http.ResponseWriter, path string) {

	if utils.IsDirectory(path) {
		path = path + "index.html"
	}

	var ok bool
	fileBytes := make([]byte, 1024)
	fileBytes, ok = staticFileCache[path]
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
	w.Write(fileBytes)
}
