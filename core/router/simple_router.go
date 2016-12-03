package router

import (
	"log"
	"regexp"
	"strings"

	"fmt"

	"github.com/fudali113/doob/utils"
)

/**
 * TODO 将储存改变，每次添加前判断是否有此url的handler存在
 * 若有则根据方法添加不同的方法handler
 */

const (

	// 各分类操作的正则表达式
	URL_CUT_SYMBOL           = "/"
	PATH_VARIABLE_SYMBOL     = "{\\w+}"
	SUFFIX_URL               = "[\\w|/]+/\\*\\*"
	PATH_VARIABLE_URL        = "[\\w|/]+{\\w+}[\\w|/]+"
	PATH_VARIABLE_SUFFIX_URL = "[\\w|/]+{\\w+}[\\w|/]+\\*\\*"
)

/**
 * url类别信息，用户对url进行分类，方便匹配的时候提高效率
 * @type {[type]}
 */
const (
	NORMAL = iota
	PATH_VARIABLE
	LAST_ALL_MATCH
	PV_AND_LAM
)

var (

	// 各分类操作的正则表达式对象
	pathVariableReg = getRegexp(PATH_VARIABLE_URL)
	lastAllMatchReg = getRegexp(SUFFIX_URL)
	pvAndLamReg     = getRegexp(PATH_VARIABLE_SUFFIX_URL)

	// 储存用户相关handler信息
	normalMap         = make(map[string]interface{}, 0)
	pathVariableMap   = make(map[int][]*pathVariableHandler, 0)
	lastAllMatchSlice = make([]*lastAllMatchHandler, 0)
)

/**
 * 储存url含有值得路由handler信息
 */
type pathVariableHandler struct {
	urlLen         int
	pathParamNames []string
	noMatchStrs    []string
	urlReg         *regexp.Regexp
	rest           interface{}
}

func (this *pathVariableHandler) getPathVariableParamMap(url string) map[string]string {
	res := make(map[string]string, 0)
	//TODO 根据url获取参数值
	resStrs := make([]string, 0)
	for i, noMatchStr := range this.noMatchStrs {
		if i == len(this.noMatchStrs)-1 || noMatchStr == "" {
			resStrs = append(resStrs, url)
			break
		}
		//log.Print("noMatchStr : ", noMatchStr)
		strs := strings.SplitN(url, noMatchStr, 2)
		//log.Print(strs)
		str := ""
		if len(strs) == 2 {
			if strs[0] == "" {
				url = strings.TrimPrefix(url, noMatchStr)
				continue
			}
			str = strs[0]
			resStrs = append(resStrs, str)
			url = strings.TrimPrefix(url, str)
		}
		url = strings.TrimPrefix(url, noMatchStr)
		//log.Print("url : ", url)
	}
	for i := 0; i < len(this.pathParamNames); i++ {
		res[this.pathParamNames[i]] = resStrs[i]
	}
	log.Print(res)
	return res
}

func (this *pathVariableHandler) String() string {
	return fmt.Sprintf("urlLen:%d,urlReg:%s", this.urlLen, this.urlReg.String())
}

/**
 * 储存尾部全匹配的相关信息
 */
type lastAllMatchHandler struct {
	prefixStr string
	rest      interface{}
}

type SimpleRouter struct {
}

func (this *SimpleRouter) Add(url string, restHandler interface{}) {
	switch getUrlClassify(url) {
	case NORMAL:
		normalMap[url] = restHandler
	case PATH_VARIABLE:
		pathVariableHandle(url, restHandler)
	case LAST_ALL_MATCH:
		lastAllMatchhandle(url, restHandler)
	case PV_AND_LAM:
		//TODO 现在暂不考虑这种情况
	}
}

func (this *SimpleRouter) Get(url string) interface{} {
	handler, ok := normalMap[url]
	if ok {
		return handler
	}
	urlStrLen := len(utils.Split(url, URL_CUT_SYMBOL))
	pvHandlers, pvOk := pathVariableMap[urlStrLen]
	if pvOk {
		for _, pvHandler := range pvHandlers {
			if pvHandler.urlReg.MatchString(url) {
				return pvHandler.rest
			}
		}
	}
	for _, lamHandler := range lastAllMatchSlice {
		if strings.Index(url, lamHandler.prefixStr) == 0 {
			return lamHandler.rest
		}
	}
	return nil
}

/**
 * 当分类为url中含有参数时的相关操作
 */
func pathVariableHandle(url string, restHandler interface{}) {
	urlStrArray := utils.Split(url, URL_CUT_SYMBOL)
	urlStrArrayLen := len(urlStrArray)
	pathVariableHandlerSlice := pathVariableMap[urlStrArrayLen]
	pathVariableHandlerSlice = append(
		pathVariableHandlerSlice,
		getPathVariableHandler(url, restHandler))
	pathVariableMap[urlStrArrayLen] = pathVariableHandlerSlice
}

func getPathVariableHandler(url string, restHandler interface{}) *pathVariableHandler {
	urlStrArray := utils.Split(url, URL_CUT_SYMBOL)
	urlStrArrayLen := len(urlStrArray)
	urlReg := getPathVariableReg(url)
	noMatchStr := getRegexp(PATH_VARIABLE_SYMBOL).Split(url, -1)
	log.Print(noMatchStr)
	return &pathVariableHandler{
		urlLen:         urlStrArrayLen,
		pathParamNames: getPathParamNames(url),
		noMatchStrs:    noMatchStr,
		urlReg:         urlReg,
		rest:           restHandler,
	}
}

func getPathParamNames(url string) []string {
	res := make([]string, 0)
	matchs := getRegexp(PATH_VARIABLE_SYMBOL).FindAllStringSubmatch(url, -1)
	for _, _match := range matchs {
		match := _match[0]
		paramName := match[1 : len(match)-1]
		res = append(res, paramName)
	}
	log.Print(url, " :paraname: ", res)
	return res
}

/**
 * 获取匹配该handler url的正则表达式
 */
func getPathVariableReg(url string) *regexp.Regexp {
	r := getRegexp(PATH_VARIABLE_SYMBOL)
	return getRegexp(r.ReplaceAllString(url, "\\S+"))
}

/**
 * 当分类为尾部全匹配时的相关操作
 */
func lastAllMatchhandle(url string, restHandler interface{}) {
	prefixStr := strings.Replace(url, "**", "", 1)
	lastAllMatchSlice = append(lastAllMatchSlice, &lastAllMatchHandler{
		prefixStr: prefixStr,
		rest:      restHandler,
	})
}

/**
 * 获取url的分类，分类处理
 */
func getUrlClassify(url string) int {
	if pvAndLamReg.MatchString(url) {
		return PV_AND_LAM
	} else if lastAllMatchReg.MatchString(url) {
		return LAST_ALL_MATCH
	} else if pathVariableReg.MatchString(url) {
		return PATH_VARIABLE
	}
	return NORMAL
}

/**
 * 获取正则表达式
 */
func getRegexp(reg string) *regexp.Regexp {
	r, err := regexp.Compile(reg)
	if err != nil {
		log.Panic(err)
		return nil
	}
	return r
}
