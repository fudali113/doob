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
	"github.com/fudali113/golib/utils"
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
	handlerMap *handleFuncMap
	filters []Filter
)

func Get_doob_Handler() *DoobHandler  {
	return &DoobHandler{filters:filters,handlerMap:handlerMap}
}

func AddHandlerFunc(url string,methodStr string, handler http.HandlerFunc){
	paras := utils.Split(url,"/")
	matchParaCount := 0
	urlinfo := &urlInfo{}
	restHandler := &restHandlerFunc{function:handler,methodStr:strings.ToLower(methodStr)}
	fmt.Println(paras,len(paras))
	for i,v := range paras{
		matchParaCount++
		para := strings.TrimSpace(v)
		//para = strings.ToLower(para)
		if para[0] == URL_PARA_PREFIX_FLAG[0] && para[len(para) - 1] == URL_PARA_LAST_FLAG[0] {
			urlinfo.addUrlPara(urlMacthPara{urlPara:para[1:len(para)-1],matchInfo:URL_PARA_FLAG})
		}else if para == ALL {
			if i == len(paras) - 1 {
				urlinfo.addUrlPara(urlMacthPara{urlPara:para,matchInfo:ALL})
				handlerMap.lastAllMatch[url[:len(url)-1]] = restHandler
			}else {
				urlinfo.addUrlPara(urlMacthPara{urlPara:para,matchInfo:URL_PARA_FLAG})
			}
		}else {
			matchParaCount--
			urlinfo.addUrlPara(urlMacthPara{urlPara:para,matchInfo:EMPTY})
		}
	}
	if matchParaCount == 0 {
		handlerMap.simple[url] = restHandler
	}else {
		urlinfo.handler = restHandler
		len := urlinfo.len()
		fmt.Println(urlinfo)
		handlerMap.rest.urls[len] = append(handlerMap.rest.urls[len],urlinfo)
	}
}

func AddFilter(f Filter)  {
	filters = append(filters,f)
}

type Filter interface {
	Filter(http.ResponseWriter,*http.Request) bool
}

type DoobHandler struct {
	filters []Filter
	handlerMap *handleFuncMap
}

func (this *DoobHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	for i,_ := range this.filters{
		if this.filters[i].Filter(res,req) {
			continue
		}else {
			return
		}
	}
	handler,err := this.handlerMap.getHandler(req)
	if err != nil{
		fmt.Println(err)
		errStr := err.Error()
		if strings.Index(errStr,"method not match") >= 0 {
			res.WriteHeader(405)
		}else{
			res.WriteHeader(404)
		}
		return
	}
	handler(res,req)
}

type handleFuncMap struct {
	simple map[string]*restHandlerFunc
	rest *restHandlerMap
	lastAllMatch map[string]*restHandlerFunc
}

func (this handleFuncMap) getHandler(req *http.Request) (http.HandlerFunc,error)  {
	url := req.URL.Path
	err := fmt.Errorf("url:%s;not found macth url",url)
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
				if req.Form == nil {
					req.Form = map[string][]string{}
				}
				req.Form.Add(k,v)
			}
			return restHandler.function,nil
		}else {
			err = fmt.Errorf("url:%s;method not match:methodStr is %s but %s",url,restHandler.methodStr,method)
		}
	}
	for k,v:=range this.lastAllMatch  {
		if index:=strings.Index(url,k);index == 0 || index == 1{
			if v.matchMethod(method){
				return v.function,nil
			}else {
				err = fmt.Errorf("url:%s;method not match:methodStr is %s but %s",url,rest.methodStr,method)
			}
		}
	}

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
	case ALL:
		return true,ALL
	case URL_PARA_FLAG:
		return true,URL_PARA_FLAG
	case EMPTY:
		return this.urlPara == para , EMPTY
	default:
		return this.urlPara == para , EMPTY
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
	urlParavalueMap := map[string]string{}
	urlParas := utils.Split(url,"/")
	for i , _ := range this.urlParas {
		should := this.urlParas[i]
		real := urlParas[i]
		if ismacth,flag:=should.macth(real); ismacth {
			switch flag {
			case ALL:
				urlParavalueMap["*"]=strings.Join(urlParas[i:],"/")
				return this.handler,urlParavalueMap,nil
			case URL_PARA_FLAG:
				urlParavalueMap[should.urlPara] = real
			default :

			}
		}else {
			return nil,nil,fmt.Errorf("url not macth")
		}
	}
	return this.handler,urlParavalueMap,nil
}

func init()  {
	simple := map[string]*restHandlerFunc{}
	rest := &restHandlerMap{urls:map[int][]*urlInfo{}}
	last := map[string]*restHandlerFunc{}
	handlerMap = &handleFuncMap{
		simple:simple,
		rest:rest,
		lastAllMatch:last,
	}
	filters = []Filter{}
}