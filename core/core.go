//
// 实际的逻辑处理
// 并对各模块进行组装最后提供给外界完整的功能
//
package core

import (
	"fmt"
	"net/http"

	"log"

	"github.com/fudali113/doob/core/router"
	"github.com/fudali113/doob/utils"
)

const (
	url_split_symbol = "&&"
)

var (
	filters = make([]Filter, 0)
	root    = router.GetRoot()

	_doob = &doob{
		filters: filters,
		root:    root,
	}
)

func Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), _doob)
}

func AddFilter(fs ...Filter) {
	_doob.addFilter(fs...)
}

// register a handler
// if urlstr use `&&` split more url
// split range register
func AddHandlerFunc(allUrl string, handler interface{}, methods ...HttpMethod) {
	urls := utils.Split(allUrl, url_split_symbol)
	for _, url := range urls {
		for _, method := range methods {
			methodStr := string(method)
			if checkMethodStr(methodStr) {
				_doob.addRestHandler(url, router.GetSimpleRestHandler(methodStr, handler))
			} else {
				log.Printf("%s method is unsupport", methodStr)
			}
		}
	}

}
