package core

import (
	"log"
	"net/http"
	"time"

	"github.com/fudali113/doob/core/errors"
)

const (
	ALL                  = "*"
	URL_PARA_PREFIX_FLAG = "{"
	URL_PARA_LAST_FLAG   = "}"
	URL_PARA_FLAG        = "{}"
	EMPTY                = ""
)

//handler与filter容器
var (
	handlerMap *handleFuncMap
	filters    []Filter
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
	handler, errs := this.handlerMap.getHandler(req)
	if len(errs) != 0 {
		for i := 0; i < len(errs); i++ {
			err := errs[i]
			switch err.(type) {
			case *errors.MethodMacthError:
				res.WriteHeader(405)
				break
			case *errors.URLMacthError:
				res.WriteHeader(404)
			}
		}
		return
	}
	if handler == nil {
		res.WriteHeader(404)
		return
	}
	handler(res, req)
}

func init() {
	simple := map[string]*restHandlerFunc{}
	rest := &restHandlerMap{urls: map[int][]*urlInfo{}}
	last := map[string]*restHandlerFunc{}
	handlerMap = &handleFuncMap{
		simple:       simple,
		rest:         rest,
		lastAllMatch: last,
	}
	filters = []Filter{}
	// AddHandlerFunc("/", "", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("doob "))
	// })
}
