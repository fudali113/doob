//
// 实际的逻辑处理
// 并对各模块进行组装最后提供给外界完整的功能
//
package doob

import (
	"fmt"
	"net/http"

	"github.com/fudali113/doob/router"
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

func DefaultRouter() Router {
	return Router{node:root}
}

func GetRouter(prefix string) Router {
	node := root.GetNode(prefix)
	return Router{node:node}
}

func Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), _doob)
}

func AddFilter(fs ...Filter) {
	_doob.addFilter(fs...)
}

