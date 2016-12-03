package router

import (
	"log"
	"regexp"
	"strings"

	"github.com/fudali113/doob/utils"
)

const (
	// 各分类操作的正则表达式
	SUFFIX_SYMBOL               = "[\\w|/]+/\\*\\*"
	PATH_VARIABLE_SYMBOL        = "[\\w|/]+{\\w+}[\\w|/]+"
	PATH_VARIABLE_SUFFIX_SYMBOL = "[\\w|/]+{\\w+}[\\w|/]+\\*\\*"
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
	pathVariableReg = getRegexp(PATH_VARIABLE_SYMBOL)
	lastAllMatchReg = getRegexp(SUFFIX_SYMBOL)
	pvAndLamReg     = getRegexp(PATH_VARIABLE_SUFFIX_SYMBOL)

	// 储存用户相关handler信息
	normalMap         = make(map[string]*interface{}, 0)
	pathVariableMap   = make(map[int][]*pathVariableHandler, 0)
	lastAllMatchSlice = make([]*lastAllMatchHandler, 0)
)

type pathVariableHandler struct {
	urlLen int
	regStr string
	rest   *interface{}
}

type lastAllMatchHandler struct {
	prefixStr string
	rest      *interface{}
}

type SimpleRouter struct {
}

func (this *SimpleRouter) Add(url string, restHandler *interface{}) {
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

func (this *SimpleRouter) Get(url string) *interface{} {
	return nil
}

/**
 * 当分类为url中含有参数时的相关操作
 */
func pathVariableHandle(url string, restHandler *interface{}) {
	urlStrArray := utils.Split(url, "/")
	urlStrArrayLen := len(urlStrArray)
	urlReg := getPathVariableReg(url)
	pathVariableHandlerSlice := pathVariableMap[urlStrArrayLen]
	pathVariableHandlerSlice = append(pathVariableHandlerSlice,
		&pathVariableHandler{
			urlLen: urlStrArrayLen,
			regStr: urlReg,
			rest:   restHandler,
		})
}

/**
 * 获取匹配该handler url的正则表达式
 */
func getPathVariableReg(url string) string {
	r := getRegexp("{\\w+}")
	return r.ReplaceAllString(url, "\\w+")
}

/**
 * 当分类为尾部全匹配时的相关操作
 */
func lastAllMatchhandle(url string, restHandler *interface{}) {
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
