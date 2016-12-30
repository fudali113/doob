package doob

import (
	"log"
	"net/http"
	"time"

	"github.com/fudali113/doob/errors"
	"github.com/fudali113/doob/router"

	. "github.com/fudali113/doob/http_const"
	mw "github.com/fudali113/doob/middleware"
)

type doob struct {
	root    *router.Node
	filters []mw.BeforeFilter
}

// 实现 http Handle 接口
func (this *doob) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer log.Printf("程序处理共消耗:%d ns", time.Now().Sub(startTime).Nanoseconds())
	// TODO user can register err deal
	defer func() {
		if err := recover(); err != nil {
			errors.CheckErr(err, w, req, isDev)
		}
	}()
	for i := range this.filters {
		if this.filters[i].DoBeforeFilter(w, req) {
			continue
		} else {
			return
		}
	}

	url := req.URL.Path
	paramMap := make(map[string]string, 0)
	handler, err := this.root.GetRT(url, paramMap)
	if err != nil {
		w.WriteHeader(NOT_FOUND)
		return
	}
	matchResult := &router.MatchResult{
		Rest:     handler,
		ParamMap: paramMap,
	}
	invoke(matchResult, w, req)

}
