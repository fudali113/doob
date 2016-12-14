package core

import (
	"net/http"
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

// 实现 http Handle 接口
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
	matchResult := this.router.Get(url)
	invoke(matchResult, w, req)

}

func (this *doob) addFilter(fs ...Filter) {
	this.filters = append(this.filters, fs...)
}

func (this *doob) addRestHandler(url string, restHandler router.RestHandler) {
	this.router.Add(url, restHandler)
}
