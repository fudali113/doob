package core

import "net/http"

/**
save handlers
*/
type restHandlerMap struct {
	urls map[int][]*urlInfo
}

/**

 */
type urlInfo struct {
	urlParas []urlMacthPara
	handler  *restHandlerFunc
}

type urlMacthPara struct {
	urlPara   string
	matchInfo string
}

type handleFuncMap struct {
	simple       map[string]*restHandlerFunc
	rest         *restHandlerMap
	lastAllMatch map[string]*restHandlerFunc
}

type restHandlerFunc struct {
	methodStr string
	function  http.HandlerFunc
}

type Filter interface {
	Filter(http.ResponseWriter, *http.Request) bool
}
