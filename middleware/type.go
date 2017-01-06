package middlerware

import "net/http"

var Middlerwares = make([]Middleware, 0)

func AddMiddlerware(fs ...Middleware) {
	Middlerwares = append(Middlerwares, fs...)
}

// Filter接口
type BeforeFilter interface {
	// Filter 的实际操作
	// 返回 bool 值决定是否通过此 filter
	DoBeforeFilter(res http.ResponseWriter, req *http.Request) bool
}

type LaterFilter interface {
	DoLaterFilter(res http.ResponseWriter, req *http.Request)
}

type Middleware interface {
	BeforeFilter
	LaterFilter
}

type HeadHTTPMethodDealer func(res http.ResponseWriter, req *http.Request)

func (this HeadHTTPMethodDealer) DoLaterFilter(res http.ResponseWriter, req *http.Request) {
	this(res, req)
}
