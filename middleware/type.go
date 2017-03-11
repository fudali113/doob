package middleware

import "net/http"

var Middlewares = make([]Middleware, 0)
var BFilters = make([]BeforeFilter, 0)
var LFilters = make([]LaterFilter, 0)

func AddMiddleware(fs ...Middleware) {
	Middlewares = append(Middlewares, fs...)
}

func AddBFilter(fs ...BeforeFilter) {
	BFilters = append(BFilters, fs...)
}

// Filter接口
type BeforeFilter interface {
	// Filter 的实际操作
	// 返回 bool 值决定是否通过此 filter
	DoBeforeFilter(res http.ResponseWriter, req *http.Request) bool
}

type BeforeFilterFunc func(res http.ResponseWriter, req *http.Request) bool

func (bff BeforeFilterFunc) DoBeforeFilter(res http.ResponseWriter, req *http.Request) bool {
	return bff(res, req)
}

type LaterFilter interface {
	DoLaterFilter(res http.ResponseWriter, req *http.Request)
}

type Middleware interface {
	BeforeFilter
	LaterFilter
}

type DefaultMiddleware struct {
}

func (DefaultMiddleware) DoBeforeFilter(res http.ResponseWriter, req *http.Request) bool {
	return true
}

func (DefaultMiddleware) DoLaterFilter(res http.ResponseWriter, req *http.Request) {
}

type HeadHTTPMethodDealer func(res http.ResponseWriter, req *http.Request)

func (this HeadHTTPMethodDealer) DoLaterFilter(res http.ResponseWriter, req *http.Request) {
	this(res, req)
}
