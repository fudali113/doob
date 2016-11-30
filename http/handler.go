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
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fudali113/golib/utils"
)

const (
	ALL                  = "*"
	URL_PARA_PREFIX_FLAG = "{"
	URL_PARA_LAST_FLAG   = "}"
	URL_PARA_FLAG        = "{}"
	EMPTY                = ""
)

//handler与filter容器
var (
	handlerMap *handleFuncMap
	filters    []Filter
)

func Start(port int) {
	log.Printf("server is starting , listen port is %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &DoobHandler{filters: filters, handlerMap: handlerMap})
	if err != nil {
		log.Printf("start is fail => %s", err.Error())
	}
}

/**
 * 注册一个handler
 */
func AddHandlerFunc(url string, methodStr string, handler http.HandlerFunc) {
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

/**
 * 添加一个过滤器
 */
func AddFilter(f Filter) {
	filters = append(filters, f)
}

/**
 * 实现http接口
 */
func (this *DoobHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer log.Printf("程序处理共消耗:%d ns", time.Now().Sub(startTime).Nanoseconds())
	for i := range this.filters {
		if this.filters[i].Filter(res, req) {
			continue
		} else {
			return
		}
	}
	handler, err := this.handlerMap.getHandler(req)
	if err != nil {
		log.Println("error => ", err.Error())
		errStr := err.Error()
		if strings.Index(errStr, "method not match") >= 0 {
			res.WriteHeader(405)
		} else {
			res.WriteHeader(404)
		}
		return
	}
	handler(res, req)
}

func init() {
	simple := map[string]*restHandlerFunc{}
	rest := &restHandlerMap{urls: map[int][]*urlInfo{}}
	last := map[string]*restHandlerFunc{}
	handlerMap = &handleFuncMap{
		simple:       simple,
		rest:         rest,
		lastAllMatch: last,
	}
	filters = []Filter{}
	AddHandlerFunc("/", "", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("doob "))
	})
}
