package core

import (
	"net/http"
	"strings"
	"time"

	"github.com/fudali113/doob/core/router"
	"github.com/fudali113/doob/log"
)

var (
	logger = log.GetLog("simple router")
)

type doob struct {
	router  router.Router
	filters []Filter
}

func (this *doob) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer logger.Info("程序处理共消耗:%d ns", time.Now().Sub(startTime).Nanoseconds())

	for i := range this.filters {
		if this.filters[i].doFilter(res, req) {
			continue
		} else {
			return
		}
	}

	url := req.URL.Path
	method := strings.ToLower(req.Method)
	matchResult := this.router.Get(url)

	if matchResult == nil {
		logger.Notice("no match url : %s", url)
		res.WriteHeader(404)
		return
	}

	handler := matchResult.Rest.GetHandler(method)
	if handler == nil {
		logger.Notice("match url : %s , but method con`t match", url)
		res.WriteHeader(405)
		return
	}

	urlParam := matchResult.ParamMap

	for k, v := range urlParam {
		if req.Form == nil {
			req.Form = map[string][]string{}
		}
		req.Form.Add(k, v)
	}

	resultHandler, _ := handler.(http.HandlerFunc)
	resultHandler(res, req)

}

func (this *doob) addFilter(fs ...Filter) {
	this.filters = append(this.filters, fs...)
}

func (this *doob) addRestHandler(url string, restHandler router.RestHandler) {
	this.router.Add(url, restHandler)
}
