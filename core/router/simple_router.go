package router

import (
	"log"
	"regexp"
	"strings"

	"fmt"

	"github.com/fudali113/doob/utils"
)

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
	urlLen int
	urlReg *regexp.Regexp
	rest   interface{}
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
	urlReg := getPathVariableReg(url)
	pathVariableHandlerSlice := pathVariableMap[urlStrArrayLen]
	pathVariableHandlerSlice = append(pathVariableHandlerSlice,
		&pathVariableHandler{
			urlLen: urlStrArrayLen,
			urlReg: urlReg,
			rest:   restHandler,
		})
	pathVariableMap[urlStrArrayLen] = pathVariableHandlerSlice
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
