package core

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/fudali113/doob/core/register"
	"github.com/fudali113/doob/core/router"
	"github.com/fudali113/doob/log"
	reflectUtils "github.com/fudali113/doob/utils/reflect"
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

	if urlParam != nil {
		for k, v := range urlParam {
			if req.Form == nil {
				req.Form = map[string][]string{}
			}
			req.Form.Add(k, v)
		}
	}

	returns := []reflect.Value{}
	if matchResult.RegisterType != nil {
		switch matchResult.RegisterType.ParamType.Type {
		case register.PARAM_NONE:
			returns = reflectUtils.Invoke(handler)
		case register.CTX:
			var contxt interface{} = getContext(res, req)
			returns = reflectUtils.Invoke(handler, contxt)
		case register.CI_PATHVARIABLE:
			urlParamValues := []interface{}{}
			for i := 0; i < matchResult.RegisterType.ParamType.CiLen; i++ {
				urlParamValues = append(urlParamValues, urlParam[matchResult.ParamNames[i]])
			}
			returns = reflectUtils.Invoke(handler, urlParamValues...)
		}
		logger.Debug("%v", returns)
		res.Write([]byte("returns"))
		return
	}

	resultHandler, _ := handler.(http.HandlerFunc)
	resultHandler(res, req)

}

func getContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		request:  req,
		response: res,
		Params:   map[string]string{},
	}
}

func (this *doob) addFilter(fs ...Filter) {
	this.filters = append(this.filters, fs...)
}

func (this *doob) addRestHandler(url string, restHandler router.RestHandler) {
	this.router.Add(url, restHandler)
}
