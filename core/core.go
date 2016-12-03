package core

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fudali113/doob/utils"
)

func Listen(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), &DoobHandler{filters: filters, handlerMap: handlerMap})
}

func AddFilter(f Filter) {
	filters = append(filters, f)
}

/**
 * 注册一个handler
 */
func AddHandlerFunc(url string, handler http.HandlerFunc, methods ...HttpMethod) {
	methodStr, err := convertHttpMethods2String(methods)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	paras := utils.Split(url, "/")
	matchParaCount := 0
	urlinfo := &urlInfo{}
	restHandler := &restHandlerFunc{function: handler, methodStr: strings.ToLower(methodStr)}
	for i, v := range paras {
		matchParaCount++
		para := strings.TrimSpace(v)
		//para = strings.ToLower(para)
		if para[0] == URL_PARA_PREFIX_FLAG[0] && para[len(para)-1] == URL_PARA_LAST_FLAG[0] {
			urlinfo.addUrlPara(urlMacthPara{urlPara: para[1 : len(para)-1], matchInfo: URL_PARA_FLAG})
		} else if para == ALL {
			if i == len(paras)-1 {
				urlinfo.addUrlPara(urlMacthPara{urlPara: para, matchInfo: ALL})
				handlerMap.lastAllMatch[url[:len(url)-1]] = restHandler
			} else {
				urlinfo.addUrlPara(urlMacthPara{urlPara: para, matchInfo: URL_PARA_FLAG})
			}
		} else {
			matchParaCount--
			urlinfo.addUrlPara(urlMacthPara{urlPara: para, matchInfo: EMPTY})
		}
	}
	if matchParaCount == 0 {
		handlerMap.simple[url] = restHandler
	} else {
		urlinfo.handler = restHandler
		len := urlinfo.len()
		handlerMap.rest.urls[len] = append(handlerMap.rest.urls[len], urlinfo)
	}
}
