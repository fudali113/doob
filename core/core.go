package core

import (
	"fmt"
	"net/http"

	"github.com/fudali113/doob/core/router"
)

var (
	_doob = &doob{
		filters: make([]Filter, 0),
		router:  &router.SimpleRouter{},
	}
)

func Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), _doob)
}

func AddFilter(fs ...Filter) {
	_doob.addFilter(fs...)
}

// 注册一个handler
func AddHandlerFunc(url string, handler interface{}, methods ...HttpMethod) {
	for _, method := range methods {
		methodStr := string(method)
		if checkMethodStr(methodStr) {
			_doob.addRestHandler(url, router.GetSimpleRestHandler(methodStr, handler))
		} else {
			logger.Notice("%s method is unsupport", methodStr)
		}
	}
}
