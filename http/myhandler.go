package http

import (
	"net/http"
	"golib/utils"
	"strings"
	"fmt"
)


var (

)

type Handler interface {
	handler(url string) bool
}

type handleFuncMap struct {
	simple map[string]restHandlerFunc
	rest *restHandlerMap
}

func (this handleFuncMap) getHandler(req *http.Request) (http.HandlerFunc,error)  {
	var err error
	url := req.URL.Path
	method := req.Method
	rest,ok:=this.simple[url]
	if ok {
		if rest.matchMethod(method) {
			return rest.function,nil
		}else {
			err = fmt.Errorf("url:%s;method not match:methodStr is %s but %s",url,rest.methodStr,method)
		}
	}
	restHandler,urlValues := this.rest.getHandler(url)
	if restHandler != nil {
		if restHandler.matchMethod(method){
			for k,v:=range urlValues{
				req.Form.Add(k,v)
			}
			return restHandler.function,nil
		}
	}
	err = fmt.Errorf("url:%s;not found macth url",url)
	return nil,err
}

type restHandlerFunc struct {
	methodStr string
	function http.HandlerFunc
}

func (this restHandlerFunc) matchMethod(method string) bool{
	if this.methodStr == "" || this.methodStr == "*" {
		return true
	}
	return strings.Index(this.methodStr, method) >= 0
}


type restHandlerMap struct {
	urls map[int][]urlInfo
}

func (this restHandlerMap) getHandler(url string) (restHandlerFunc,map[string]string) {
	urlParaLen := len(utils.Split(url,"/"))
	urlInfos,ok := this.urls[urlParaLen]
	if ok {
		for _,v := range urlInfos {
			handler,urlPara,err := v.match(url)
			if err != nil{
				continue
			}else{
				return handler,urlPara
			}
		}
	}
	return nil,nil
}

func (this Handler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	url := req.URL.Path
	urlParas := utils.Split(url , "/")
	urlLenGroup := urls[len(urlParas)]

}

type urlMacthPara struct {
	urlPara string
	matchInfo string
}

func (this urlMacthPara) macth(para string) (bool,string) {
	switch this.matchInfo {
	case "*":
		return true,"*"
	case "{}":
		return true,para
	case "":
		return this.urlPara == para , ""
	default:
		return this.urlPara == para , ""
	}
}

type urlInfo struct {
	urlParas []urlMacthPara
	handler restHandlerFunc
}

func (this urlInfo) addUrlPara(para string)  {
	this.urls = append(this.urls , para)
}

func (this urlInfo) addUrlMatchPara(para string)  {
	matchIndex := len(this.urls)
	this.matchIndex = append(this.matchIndex , matchIndex)
	this.urls = append(this.urls , para)
}

func (this urlInfo) len() int  {
	return len(this.urls)
}

func (this urlInfo) match(url string) (http.HandlerFunc,map[string]string,error) {
	var urlParavalueMap map[string]string
	urlParas := utils.Split(url,"/")
	for i , _ := range urlParas {
		should := this.urlParas[i]
		real := urlParas[i]
		if macth,res:=should.macth(real);macth{
			switch res {
			case "*":
				urlParavalueMap["*"]=strings.Join(urlParas[i:],"/")
				return this.handler,urlParavalueMap,nil
			case "{}":
				urlParavalueMap[should.urlPara] = real
			default :

			}
		}else {
			return nil,nil,fmt.Errorf("not macth")
		}
	}
	return this.handler,urlParavalueMap,nil
}

func AddHandlerFunc(url string, handler http.HandlerFunc){
	paras := utils.Split(url,"/")
	for _,v := range paras{
		urlinfo := urlInfo{}
		para := strings.TrimSpace(v)
		if para[0] == '{' && para[len(para) - 1] == '}' {
			urlinfo.addUrlMatchPara(para[1:len(para)])
		}else if para == '*' {
			urlinfo.addUrlMatchPara(para)
		}else {
			urlinfo.addUrlPara(para)
			simpleUrls[url] = handler
			return
		}
		urlinfo.handler = handler
		len := urlinfo.len()
		urls[len] = append(urls[len],urlinfo)
	}
}