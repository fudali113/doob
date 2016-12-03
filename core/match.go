package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fudali113/doob/core/errors"
	"github.com/fudali113/doob/utils"
)

/**
 * 根据request获取相应的处理handler
 */
func (this handleFuncMap) getHandler(req *http.Request) (http.HandlerFunc, []error) {
	url := req.URL.Path
	errs := []error{}
	method := strings.ToLower(req.Method)
	rest, ok := this.simple[url]
	if ok {
		if rest.matchMethod(method) {
			return rest.function, nil
		}
		errs = append(errs, errors.GetMwthodMatchError(rest.methodStr, method, url, nil))
	}
	restHandler, urlValues := this.rest.getHandler(url)
	if restHandler != nil {
		if restHandler.matchMethod(method) {
			for k, v := range urlValues {
				if req.Form == nil {
					req.Form = map[string][]string{}
				}
				req.Form.Add(k, v)
			}
			return restHandler.function, nil
		}
		errs = append(errs, errors.GetMwthodMatchError(rest.methodStr, method, url, nil))
	}
	for k, v := range this.lastAllMatch {
		if index := strings.Index(url, k); index == 0 || index == 1 {
			if v.matchMethod(method) {
				return v.function, nil
			}
			errs = append(errs, errors.GetMwthodMatchError(rest.methodStr, method, url, nil))
		}
	}

	return nil, errs
}

/**
 * 根据路由匹配仓库中的handler
 */
func (this *urlInfo) match(url string) (*restHandlerFunc, map[string]string, error) {
	urlParavalueMap := map[string]string{}
	urlParas := utils.Split(url, "/")
	for i, _ := range this.urlParas {
		should := this.urlParas[i]
		real := urlParas[i]
		if ismacth, flag := should.macth(real); ismacth {
			switch flag {
			case ALL:
				urlParavalueMap["*"] = strings.Join(urlParas[i:], "/")
				return this.handler, urlParavalueMap, nil
			case URL_PARA_FLAG:
				urlParavalueMap[should.urlPara] = real
			default:

			}
		} else {
			return nil, nil, fmt.Errorf("url not macth")
		}
	}
	return this.handler, urlParavalueMap, nil
}

func (this restHandlerMap) getHandler(url string) (*restHandlerFunc, map[string]string) {
	urlParaLen := len(utils.Split(url, "/"))
	urlInfos, ok := this.urls[urlParaLen]
	if ok {
		for _, v := range urlInfos {
			handler, urlPara, err := v.match(url)
			if err != nil {
				continue
			} else {
				return handler, urlPara
			}
		}
	}
	return nil, nil
}

func (this urlMacthPara) macth(para string) (bool, string) {
	switch this.matchInfo {
	case ALL:
		return true, ALL
	case URL_PARA_FLAG:
		return true, URL_PARA_FLAG
	case EMPTY:
		return this.urlPara == para, EMPTY
	default:
		return this.urlPara == para, EMPTY
	}
}

/**
 * 验证http 方法是否匹配
 */
func (this restHandlerFunc) matchMethod(method string) bool {
	method = strings.ToLower(method)
	methodIsMatch := this.methodStr == "" || this.methodStr == "*"
	if !methodIsMatch {
		return strings.Index(this.methodStr, method) >= 0
	}
	return true
}
