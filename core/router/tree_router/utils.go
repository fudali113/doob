package tree_router

import (
	"log"
	"strings"

	"github.com/fudali113/doob/utils"
)

const (
	urlSplitSymbol = "/"
	pathVarRegStr  = "{\\S\\s+}"
)

var (
	pathVarReg = utils.GetRegexp(pathVarRegStr)
)

func getUrlNodeValue(url string) (string, string) {
	url = strings.TrimPrefix(url, "/")
	prefixAndSuffix := strings.SplitN(url, "/", 1)
	return prefixAndSuffix[0], prefixAndSuffix[1]
}

// 获取url的前缀和剩余的部分
func splitUrl(url string) (string, string) {
	url = strings.TrimPrefix(url, urlSplitSymbol)
	prefixAndOther := strings.SplitN(url, urlSplitSymbol, 2)
	if len(prefixAndOther) == 1 {
		return prefixAndOther[0], ""
	}
	return prefixAndOther[0], prefixAndOther[1]
}

// 根据url的每一部分生成一个特定的value
// 用于寻找路径是做匹配
func createNodeValue(urlPart string) nodeV {
	if strings.HasSuffix(urlPart, "**") {
		prefixStr := strings.Replace(urlPart, "**", "", 1)
		if strings.HasSuffix(prefixStr, "**") {
			log.Panic("只允许在最后添加**")
		}
		return &nodeVMatchAll{
			origin: urlPart,
			prefix: prefixStr,
		}
	}

	if matchStr := pathVarReg.FindAllString(urlPart, -1); len(matchStr) > 0 {
		pathVarStr := matchStr[0]
		if paramNameAndReg := utils.Split(pathVarStr, ":"); len(paramNameAndReg) > 1 {
			parLen := len(paramNameAndReg)
			return &nodeVPathReg{
				origin:    urlPart,
				paramName: paramNameAndReg[0],
				paramReg:  utils.GetRegexp(strings.Join(paramNameAndReg[1:parLen-1], "")),
			}
		}
		return &nodeVPathVar{
			origin:    urlPart,
			paramName: pathVarStr,
		}
	}
	return &nodeVNormal{origin: urlPart}
}

func isMatch(fact, origin string, class int) bool {
	switch class {
	case normal:
		return fact == origin
	case pathReg:
		return false
	}
	return false
}

// 根据参数获取参数类别
func getClass(s string) int {
	if s == "**" {
		return matchAll
	}
	if matchStr := pathVarReg.FindAllString(s, -1); len(matchStr) > 0 {
		pathVarStr := matchStr[0]
		if strings.Contains(pathVarStr, ":") {
			return pathReg
		}
		return pathVar
	}
	return normal
}
