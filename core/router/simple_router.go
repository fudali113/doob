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
	emptyMap = make(map[string]string, 0)

	// 各分类操作的正则表达式对象
	pathVariableReg = getRegexp(PATH_VARIABLE_URL)
	lastAllMatchReg = getRegexp(SUFFIX_URL)
	pvAndLamReg     = getRegexp(PATH_VARIABLE_SUFFIX_URL)

	// 储存用户相关handler信息
	normalMap         = make(map[string]RestHandler, 0)
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
	rest           RestHandler
}

/**
 * 根据实际url获取path variable
 */
func (this *pathVariableHandler) getPathVariableParamMap(url string) map[string]string {
	res := make(map[string]string, 0)
	resStrs := make([]string, 0)
	for i, noMatchStr := range this.noMatchStrs {
		/**
		 * 如果是最后一个值了且分割字符串为空
		 * 则代表最后一个字符串为想要获取的字符串
		 */
		if i == len(this.noMatchStrs)-1 || noMatchStr == "" {
			resStrs = append(resStrs, url)
			break
		}
		/**
		 * 将数组按分割字符串分为两组
		 * 当第一个为空时，说明之前没有想要获取的字符串
		 * 当不为空时，说明包含想要获取的字符串，获取这个字符串
		 * 并将获取的字符串和分割字符串从原始字符串中去掉
		 */
		strs := strings.SplitN(url, noMatchStr, 2)
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
	rest      RestHandler
}

/**
 * 获取尾部匹配的string
 */
func (this *lastAllMatchHandler) getLastStr(url string) string {
	return strings.TrimRight(url, this.prefixStr)
}

/**
 * 返回值类型
 */
type GetResult struct {
	ParamMap map[string]string
	Rest     RestHandler
}

type SimpleRouter struct {
}

func (this *SimpleRouter) Add(url string, restHandler RestHandler) {
	switch getUrlClassify(url) {
	case NORMAL:
		_restHandler, ok := normalMap[url]
		if ok {
			_restHandler.Joint(restHandler)
			normalMap[url] = _restHandler
		} else {
			normalMap[url] = restHandler
		}
	case PATH_VARIABLE:
		pathVariableHandle(url, restHandler)
	case LAST_ALL_MATCH:
		lastAllMatchhandle(url, restHandler)
	case PV_AND_LAM:
		//TODO 现在暂不考虑这种情况
	}
}

func (this *SimpleRouter) Get(url string) *GetResult {
	handler, ok := normalMap[url]
	if ok {
		return &GetResult{
			ParamMap: emptyMap,
			Rest:     handler,
		}
	}
	urlStrLen := len(utils.Split(url, URL_CUT_SYMBOL))
	pvHandlers, pvOk := pathVariableMap[urlStrLen]
	if pvOk {
		for _, pvHandler := range pvHandlers {
			if pvHandler.urlReg.MatchString(url) {
				return &GetResult{
					Rest:     pvHandler.rest,
					ParamMap: pvHandler.getPathVariableParamMap(url),
				}
			}
		}
	}
	for _, lamHandler := range lastAllMatchSlice {
		if strings.Index(url, lamHandler.prefixStr) == 0 {
			return &GetResult{
				Rest: lamHandler.rest,
				ParamMap: map[string]string{
					"last": lamHandler.getLastStr(url),
				},
			}
		}
	}
	return nil
}

/**
 * 当分类为url中含有参数时的相关操作
 */
func pathVariableHandle(url string, restHandler RestHandler) {
	urlStrArray := utils.Split(url, URL_CUT_SYMBOL)
	urlStrArrayLen := len(urlStrArray)
	pathVariableHandlerSlice := pathVariableMap[urlStrArrayLen]
	for _, pvHandler := range pathVariableHandlerSlice {
		if pvHandler.urlReg.String() == getPathVariableReg(url).String() {
			pvHandler.rest.Joint(restHandler)
			return
		}
	}
	pathVariableHandlerSlice = append(
		pathVariableHandlerSlice,
		getPathVariableHandler(url, restHandler))
	pathVariableMap[urlStrArrayLen] = pathVariableHandlerSlice
}

/**
 * 根据用户注册的 url 和 handler 生成一个 PathVariableHandler
 */
func getPathVariableHandler(url string, restHandler RestHandler) *pathVariableHandler {
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

/**
 * 根据url获取用户注册url中的参数名字
 */
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
func lastAllMatchhandle(url string, restHandler RestHandler) {
	prefixStr := strings.Replace(url, "**", "", 1)
	for _, lastAllMatchHandler := range lastAllMatchSlice {
		if lastAllMatchHandler.prefixStr == prefixStr {
			lastAllMatchHandler.rest.Joint(restHandler)
			return
		}
	}
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
