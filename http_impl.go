package doob

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fudali113/doob/router"

	. "github.com/fudali113/doob/http_const"
)

type doob struct {
	root    *router.Node
	filters []Filter
}

// 实现 http Handle 接口
func (this *doob) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer log.Printf("程序处理共消耗:%d ns", time.Now().Sub(startTime).Nanoseconds())
	// TODO user can register err deal
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			default:
				w.WriteHeader(INTERNAL_SERVER_ERROR)
				w.Write([]byte(fmt.Sprintf("%v", err)))
			}
		}
	}()
	for i := range this.filters {
		if this.filters[i].doFilter(w, req) {
			continue
		} else {
			return
		}
	}

	url := req.URL.Path
	paramMap := make(map[string]string, 0)
	handler, err := this.root.GetRT(url, paramMap)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	matchResult := &router.MatchResult{
		Rest:     handler,
		ParamMap: paramMap,
	}
	invoke(matchResult, w, req)

}
