package router

import (
	"log"
	"strings"

	"github.com/fudali113/doob/utils"
)

const (
	urlSplitSymbol    = "/"
	pathVarRegStr     = "{\\S+}"
	suffixMatchSymbol = "*"
)

var (
	pathVarReg = utils.GetRegexp(pathVarRegStr)
)

// splitURL 获取url的前缀和剩余的部分
func splitURL(URL string) (string, string) {
	URL = strings.TrimPrefix(URL, urlSplitSymbol)
	prefixAndOther := strings.SplitN(URL, urlSplitSymbol, 2)
	if len(prefixAndOther) == 1 {
		return prefixAndOther[0], ""
	}
	return prefixAndOther[0], prefixAndOther[1]
}

// 根据url的每一部分生成一个特定的value
// 用于寻找路径是做匹配
// bug  当初先`/d/s*`情况时将出现bug
func createNodeValue(urlPart string) nodeV {
	if strings.HasSuffix(urlPart, suffixMatchSymbol) {
		prefixStr := strings.Replace(urlPart, suffixMatchSymbol, "", 1)
		if strings.HasSuffix(prefixStr, suffixMatchSymbol) {
			log.Panic(`
			"*" 字符只允许出现在路由的最末端
			right: 
				- a/b/cap/*
			error:
				- /a/d/**
				- /*/d/c  这样你后面的参数将无效
			`)
		}
		return &nodeVMatchAll{
			origin: urlPart,
			prefix: prefixStr,
		}
	}

	if matchStr := pathVarReg.FindAllString(urlPart, -1); len(matchStr) > 0 {
		pathVarStr := strings.TrimPrefix(matchStr[0], "{")
		pathVarStr = strings.TrimSuffix(pathVarStr, "}")
		if paramNameAndReg := utils.Split(pathVarStr, ":"); len(paramNameAndReg) > 1 {
			paraLen := len(paramNameAndReg)
			paraRegStr := strings.Join(paramNameAndReg[1:paraLen], "")
			return &nodeVPathReg{
				origin:    urlPart,
				paramName: paramNameAndReg[0],
				paramReg:  utils.GetRegexp(paraRegStr),
			}
		}
		return &nodeVPathVar{
			origin:    urlPart,
			paramName: pathVarStr,
		}
	}
	return &nodeVNormal{origin: urlPart}
}

// getClass 根据参数获取参数类别
func getClass(s string) int {
	if s == suffixMatchSymbol {
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

// addValueToPathParam 将参数值添加到map中
func addValueToPathParam(paramMap map[string]string, k, v string) {
	if paramMap != nil {
		paramMap[k] = v
	}
}
