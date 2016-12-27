package doob

import "net/http"

// Filter接口
type Filter interface {
	// Filter 的实际操作
	// 返回 bool 值决定是否通过此 filter
	doFilter(res http.ResponseWriter, req *http.Request) bool
}
