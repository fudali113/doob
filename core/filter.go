package core

import "net/http"

type Filter interface {
	doFilter(res http.ResponseWriter, req *http.Request) bool
}
