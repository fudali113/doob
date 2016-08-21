package http

/*
现实思路:
    分为三种类型url
	1.普通的:如/fff/ddd/lll
	2.有映射值的:如/user/{who}/info
	3.尾部全部匹配的:如/user/*('*'只可以用于尾部)
    先进行分组,将普通的于要进行取值的分开
    普通的对于一个map,直接使用map[string]获取handlerFunc
    要取值得将url利用"/"切分成数组,匹配是再将实际url切分成数组进行比对并取值

 */

import (
	"net/http"
	"golib/utils"
	"strings"
	"fmt"
)

const (
	ALL = "*"
	URL_PARA_PREFIX_FLAG = "{"
	URL_PARA_LAST_FLAG = "}"
	URL_PARA_FLAG = "{}"
	EMPTY = ""
)

var (
	handlerMap handleFuncMap
)

type DoobHandler struct {

}

func (this DoobHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	handler,err := handlerMap.getHandler(req)
	if err != nil{
		fmt.Println(err)
		res.WriteHeader(404)
	}
	handler(res,req)
}

type handleFuncMap struct {
	simple map[string]*restHandlerFunc
	rest *restHandlerMap
}

func (this handleFuncMap) getHandler(req *http.Request) (http.HandlerFunc,error)  {
	var err error
	url := req.URL.Path
	method := strings.ToLower(req.Method)
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

func (this restHandlerFunc) matchMethod(method string) bool {
	method = strings.ToLower(method)
	if this.methodStr == "" || this.methodStr == "*" {
		return true
	}
	return strings.Index(this.methodStr, method) >= 0
}


type restHandlerMap struct {
	urls map[int][]*urlInfo
}

func (this restHandlerMap) getHandler(url string) (*restHandlerFunc,map[string]string) {
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
	handler *restHandlerFunc
}

func (this *urlInfo) addUrlPara(v urlMacthPara)  {
	this.urlParas = append(this.urlParas , v)
}

func (this *urlInfo) len() int  {
	return len(this.urlParas)
}

func (this *urlInfo) match(url string) (*restHandlerFunc,map[string]string,error) {
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

func AddHandlerFunc(url string,methodStr string, handler http.HandlerFunc){
	paras := utils.Split(url,"/")
	matchParaCount := 0
	urlinfo := &urlInfo{}
	fmt.Println(paras,len(paras))
	for i,v := range paras{
		matchParaCount++
		para := strings.TrimSpace(v)
		//para = strings.ToLower(para)
		if para[0] == URL_PARA_PREFIX_FLAG[0] && para[len(para) - 1] == URL_PARA_LAST_FLAG[0] {
			urlinfo.addUrlPara(urlMacthPara{urlPara:para[1:len(para)],matchInfo:URL_PARA_FLAG})
		}else if para == ALL {
			if i == len(paras) - 1 {
				urlinfo.addUrlPara(urlMacthPara{urlPara:para,matchInfo:ALL})
			}else {
				urlinfo.addUrlPara(urlMacthPara{urlPara:para,matchInfo:URL_PARA_FLAG})
			}
		}else {
			matchParaCount--
			urlinfo.addUrlPara(urlMacthPara{urlPara:para,matchInfo:EMPTY})
		}
	}
	restHandler := &restHandlerFunc{function:handler,methodStr:strings.ToLower(methodStr)}
	if matchParaCount == 0 {
		handlerMap.simple[url] = restHandler
	}else {
		urlinfo.handler = restHandler
		len := urlinfo.len()
		fmt.Println(urlinfo)
		handlerMap.rest.urls[len] = append(handlerMap.rest.urls[len],urlinfo)
	}

}

func init()  {
	simple := map[string]*restHandlerFunc{}
	rest := &restHandlerMap{urls:map[int][]*urlInfo{}}
	handlerMap = handleFuncMap{
		simple:simple,
		rest:rest,
	}
}