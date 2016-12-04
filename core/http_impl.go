package core

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fudali113/doob/core/router"
)

type doob struct {
	router  router.Router
	filters []Filter
}

func (this *doob) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer log.Printf("程序处理共消耗:%d ns", time.Now().Sub(startTime).Nanoseconds())

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
		log.Print("no match url")
		return
	}

	handler := matchResult.Rest.GetHandler(method)
	if handler == nil {
		log.Print("match url , but method con`t match")
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
