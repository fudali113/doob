package router

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fudali113/doob/core/register"
	"github.com/fudali113/doob/log"
	"github.com/fudali113/doob/utils"
)

/**
 * TODO 给每个url模板打个分判断权重值，最后对有路径参数的url模板数组排序
 * 若有则根据方法添加不同的方法handler
 */

const (

	// 各分类操作的正则表达式
	ALL_MATCH_REG            = "\\S+"
	URL_CUT_SYMBOL           = "/"
	PATH_VARIABLE_SYMBOL     = "{\\S+?}"
	SUFFIX_URL               = "[\\w|/]+\\*\\*"
	PATH_VARIABLE_URL        = "[\\w|/]+{\\S+}[\\w|/]*"
	PATH_VARIABLE_SUFFIX_URL = "[\\w|/]+{\\S+}[\\w|/]+\\*\\*"
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
	logger   = log.GetLog("simple router")
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
	splitStrs      []string
	urlReg         *regexp.Regexp
	rest           RestHandler
	registerType   *register.RegisterType
}

/**
 * 根据实际url获取path variable
 */
func (this *pathVariableHandler) getPathVariableParamMap(url string) map[string]string {
	res := make(map[string]string, 0)
	resStrs := make([]string, 0)
	for i, splitStr := range this.splitStrs {

		/**
		 * 为了支持在用户url模板中使用正则时可以使用一个{}符号加入的判断条件
		 * 其他无任何意义
		 */
		splitStr = strings.TrimPrefix(splitStr, "}")

		/**
		 * 如果是最后一个值了且分割字符串为空
		 * 则代表最后一个字符串为想要获取的字符串
		 */
		if i == len(this.splitStrs)-1 && splitStr == "" {
			resStrs = append(resStrs, url)
			break
		}
		/**
		 * 将数组按分割字符串分为两组
		 * 当第一个为空时，说明之前没有想要获取的字符串
		 * 当不为空时，说明包含想要获取的字符串，获取这个字符串
		 * 并将获取的字符串和分割字符串从原始字符串中去掉
		 */
		strs := strings.SplitN(url, splitStr, 2)
		str := ""
		if len(strs) == 2 {
			if strs[0] == "" {
				url = strings.TrimPrefix(url, splitStr)
				continue
			}
			str = strs[0]
			resStrs = append(resStrs, str)
			url = strings.TrimPrefix(url, str)
		}
		url = strings.TrimPrefix(url, splitStr)
	}
	for i := 0; i < len(this.pathParamNames); i++ {
		res[this.pathParamNames[i]] = resStrs[i]
	}
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
		logger.Error("url : %s %s", url, " is not add")
		//TODO 现在暂不考虑这种情况
	}
}

func (this *SimpleRouter) Get(url string) *MatchResult {
	handler, ok := normalMap[url]
	if ok {
		return &MatchResult{
			ParamMap: emptyMap,
			Rest:     handler,
		}
	}
	urlStrLen := len(utils.Split(url, URL_CUT_SYMBOL))
	pvHandlers, pvOk := pathVariableMap[urlStrLen]
	if pvOk {
		for _, pvHandler := range pvHandlers {
			if pvHandler.urlReg.MatchString(url) {
				return &MatchResult{
					Rest:         pvHandler.rest,
					ParamNames:   pvHandler.pathParamNames,
					ParamMap:     pvHandler.getPathVariableParamMap(url),
					RegisterType: pvHandler.registerType,
				}
			}
		}
	}
	for _, lamHandler := range lastAllMatchSlice {
		if strings.Index(url, lamHandler.prefixStr) == 0 {
			return &MatchResult{
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
		if pathRegexp, _ := getPathVariableRegAndParamNames(url); pvHandler.urlReg.String() == pathRegexp.String() {
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
	urlReg, pathParamNames := getPathVariableRegAndParamNames(url)
	noMatchStr := getRegexp(PATH_VARIABLE_SYMBOL).Split(url, -1)
	return &pathVariableHandler{
		urlLen:         urlStrArrayLen,
		pathParamNames: pathParamNames,
		splitStrs:      noMatchStr,
		urlReg:         urlReg,
		rest:           restHandler,
		registerType:   register.GetFuncRegisterType(restHandler.GetSigninHandler()),
	}
}

/**
 * 获取匹配该handler url的正则表达式和获取参数名
 */
func getPathVariableRegAndParamNames(url string) (*regexp.Regexp, []string) {
	r := getRegexp(PATH_VARIABLE_SYMBOL)
	paramNames := make([]string, 0)
	templates := r.FindAllString(url, -1)
	for _, template := range templates {
		templateCut := template[1 : len(template)-1]
		paramName, regexpStr := getTemplateNameAndRegexpStr(templateCut)
		paramNames = append(paramNames, paramName)
		url = strings.Replace(url, template, regexpStr, 1)
	}

	return getRegexp(url), paramNames
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
		logger.Error(err.Error())
		return nil
	}
	return r
}

/**
 * 获取用户的正则表达式
 * 如果没有则匹配全部
 */
func getTemplateNameAndRegexpStr(template string) (string, string) {
	templateArray := utils.Split(template, ":")
	if len(templateArray) == 2 {
		return templateArray[0], templateArray[1]
	}
	return templateArray[0], ALL_MATCH_REG
}
