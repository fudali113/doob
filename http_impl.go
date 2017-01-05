package doob

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fudali113/doob/errors"
	"github.com/fudali113/doob/register"
	"github.com/fudali113/doob/router"

	. "github.com/fudali113/doob/http_const"
	mw "github.com/fudali113/doob/middleware"
	returnDeal "github.com/fudali113/doob/return_deal"
	reflectUtils "github.com/fudali113/doob/utils/reflect"
)

type doob struct {
	root         *router.Node
	bFilters     []mw.BeforeFilter
	lFilters     []mw.LaterFilter
	middlerwares []mw.Middleware
}

// 实现 http Handle 接口
func (this *doob) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	url := req.URL.Path
	method := strings.ToLower(req.Method)

	defer log.Printf("url: %s | method: %s | deal time: %d ns", url, method, time.Now().Sub(startTime).Nanoseconds())
	defer func() {
		if err := recover(); err != nil {
			errors.CheckErr(err, w, req, IsDev)
		}
	}()

	// 前处理
	for i := range this.middlerwares {
		if this.bFilters[i].DoBeforeFilter(w, req) {
			continue
		} else {
			return
		}
	}

	for i := range this.bFilters {
		if this.bFilters[i].DoBeforeFilter(w, req) {
			continue
		} else {
			return
		}
	}

	paramMap := make(map[string]string, 0)
	handler, err := this.root.GetRT(url, paramMap)
	if err != nil {
		w.WriteHeader(NOT_FOUND)
		return
	}

	handlerType := handler.GetHandler(method)
	if handlerType == nil {

		if AutoAddOptions && method == string(OPTIONS) {
			methods := handler.GetMethods()
			returnDeal.DealReturn(&returnDeal.ReturnType{
				TypeStr: returnDeal.DEFAULT_JSON_DEALER_NAME,
				Data:    methods,
			}, w)
		} else {
			log.Printf("match url : %s , but method con`t match", url)
			w.WriteHeader(METHOD_NOT_ALLOWED)
		}
		return
	}

	matchResult := &router.MatchResult{
		Rest:     handler,
		ParamMap: paramMap,
	}
	invoke(matchResult, handlerType, w, req)

	// 后处理
	for i := range this.lFilters {
		this.lFilters[i].DoLaterFilter(w, req)
	}

	for i := range this.middlerwares {
		this.bFilters[len(this.middlerwares)-1-i].DoBeforeFilter(w, req)
	}
}

//
// 根据路由匹配获取匹配的返回值
// 根据返回值执行不同的逻辑操作
//
// FIXME 此方法有些复杂，需要进行拆解
//
func invoke(matchResult *router.MatchResult, handlerType register.RegisterHandlerType, w http.ResponseWriter, req *http.Request) {

	// 获取路劲参数并存入request参数中
	urlParam := matchResult.ParamMap
	if urlParam != nil {
		for k, v := range urlParam {
			if req.Form == nil {
				req.Form = map[string][]string{}
			}
			req.Form.Add(k, v)
		}
	}

	// 根据RegisterType决定怎么执行函数
	registerType := handlerType.GetRegisterType()
	handlerInterface := handlerType.GetHandler()
	if registerType != nil {
		paramType := registerType.ParamType
		returnType := registerType.ReturnType
		switch paramType.Type {

		case register.ORIGIN:
			handler := handlerInterface.(func(http.ResponseWriter, *http.Request))
			handler(w, req)

		case register.PARAM_NONE:
			switch returnType.Type {
			case register.RETURN_NONE:
				handler := handlerInterface.(func())
				handler()

			case register.JSON:
				handler := handlerInterface.(func() interface{})
				returnValue := handler()
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: ReturnDealDefaultType,
					Data:    returnValue,
				}, w)

			case register.FILE:
				handler := handlerInterface.(func() string)
				returnValue := handler()
				returnDeal.DealReturn(&returnDeal.ReturnType{TypeStr: returnValue}, w)

			case register.RETURN_TYPE:
				handler := handlerInterface.(func() (string, interface{}))
				str, data := handler()
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: str,
					Data:    data,
				}, w)
			}

		case register.CTX:
			context := getContext(w, req)
			switch returnType.Type {
			case register.RETURN_NONE:
				handler := handlerInterface.(func(*Context))
				handler(context)

			case register.FILE:
				handler := handlerInterface.(func(*Context) string)
				returnValue := handler(context)
				returnDeal.DealReturn(&returnDeal.ReturnType{TypeStr: returnValue}, w)

			case register.JSON:
				handler := handlerInterface.(func(*Context) interface{})
				returnValue := handler(context)
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: getReqAccept(req),
					Data:    returnValue,
				}, w)

			case register.RETURN_TYPE:
				handler := handlerInterface.(func(*Context) (string, interface{}))
				str, data := handler(context)
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: str,
					Data:    data,
				}, w)
			}

		case register.CI_PATHVARIABLE, register.CI_PATHVARIABLE_CTX:
			var returns []interface{}
			ciLen := paramType.CiLen
			paraNames := make([]string, 0)
			for k, _ := range matchResult.ParamMap {
				paraNames = append(paraNames, k)
			}
			if ciLen > len(paraNames) {
				log.Printf(`your func path variable params lnegth is %d ,
           but your url params length just %d`, ciLen, len(paraNames))
				return
			}

			params := []interface{}{}
			for i := 0; i < ciLen; i++ {
				params = append(params, urlParam[paraNames[i]])
			}
			if paramType.Type == register.CI_PATHVARIABLE_CTX {
				params = append(params, getContext(w, req))
			}
			returns = reflectUtils.Invoke(handlerInterface, params...)

			switch returnType.Type {
			case register.RETURN_NONE:

			case register.FILE:
				str := returns[0].(string)
				returnDeal.DealReturn(&returnDeal.ReturnType{TypeStr: str}, w)

			case register.JSON:
				returnValue := returns[0].(interface{})
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: "json",
					Data:    returnValue,
				}, w)

			case register.RETURN_TYPE:
				str := returns[0].(string)
				returnValue := returns[1]
				returnDeal.DealReturn(&returnDeal.ReturnType{
					TypeStr: str,
					Data:    returnValue,
				}, w)
			}
			return
		}
	}
}

// According to the request and response for context
func getContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Request:    req,
		Response:   w,
		PathParams: map[string]string{},
	}
}

// if user don`t set ReturnDealDefaultType
// ReturnDealDefaultType deafault value is "auto"
// Will automatically think return type according to the request to accept
func getReqAccept(req *http.Request) string {
	if ReturnDealDefaultType != "auto" {
		return ReturnDealDefaultType
	}
	accept := req.Header.Get(ACCEPT)
	if strings.Contains(accept, APP_JSON) {
		return returnDeal.JSON_DEAL_TYPE_STR
	}
	if strings.Contains(accept, APP_XML) {
		return returnDeal.XML_DEAL_TYPE_STR
	}
	return returnDeal.JSON_DEAL_TYPE_STR
}
