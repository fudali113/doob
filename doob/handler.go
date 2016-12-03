package doob

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	ALL                  = "*"
	URL_PARA_PREFIX_FLAG = "{"
	URL_PARA_LAST_FLAG   = "}"
	URL_PARA_FLAG        = "{}"
	EMPTY                = ""
)

/**
 * 实现go http handler ，接管路由的分发
 */
type DoobHandler struct {
	filters    []Filter
	handlerMap *handleFuncMap
}

/**
 * 实现http接口
 */
func (this *DoobHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer log.Printf("程序处理共消耗:%d ns", time.Now().Sub(startTime).Nanoseconds())
	for i := range this.filters {
		if this.filters[i].Filter(res, req) {
			continue
		} else {
			return
		}
	}
	handler, err := this.handlerMap.getHandler(req)
	if err != nil {
		log.Println("error => ", err.Error())
		errStr := err.Error()
		if strings.Index(errStr, "method not match") >= 0 {
			res.WriteHeader(405)
		} else {
			res.WriteHeader(404)
		}
		return
	}
	handler(res, req)
}
