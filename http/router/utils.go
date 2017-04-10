package router

import (
	"strings"

	"github.com/fudali113/doob/utils"
)

const (
	urlSplitSymbol    = "/"
	pathVarRegStr     = "{\\S+}"
	allMatchStr       = "\\S+"
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
			panic(`
			"*" 字符只允许出现在路由的最末端
			right: 
				- a/b/cap/*
			error:
				- /a/d/**
				- /*/d/c  这样你后面的参数将无效
			`)
		}
		return &nodeVMatchAll{
			Origin: urlPart,
			prefix: prefixStr,
		}
	}

	otherIsEnpty := func(others []string) bool {
		for i := 0; i < len(others); i++ {
			if others[i] != "" {
				return false
			}
		}
		return true
	}

	// 处理匹配相关value
	matchs, others := getMatchsAndOtehrs(urlPart)
	matchsLen := len(matchs)
	switch {
	case matchsLen == 1 && otherIsEnpty(others):
		name, regStr := getNameAndRegstr(matchs[0])
		if regStr == "" {
			return &nodeVPathVar{
				Origin:    urlPart,
				paramName: name,
			}
		}
		return &nodeVPathReg{
			Origin:    urlPart,
			paramName: name,
			ParamReg:  utils.GetRegexp(regStr),
		}
	default:
		finallyRegStr := ""
		names := []string{}
		regs := []string{}
		for _, v := range matchs {
			name, reg := getNameAndRegstr(v)
			if reg == "" {
				reg = allMatchStr
			}
			names = append(names, name)
			regs = append(regs, reg)
		}
		for i := 0; i < len(regs); i++ {
			finallyRegStr += others[i]
			finallyRegStr += regs[i]
		}
		return &nodeVPathReg{
			Origin:    urlPart,
			paramName: strings.Join(names, ","),
			ParamReg:  utils.GetRegexp(finallyRegStr),
		}
	}
	return &nodeVNormal{Origin: urlPart}
}

// getNameAndRegstr matchStr需要以`{`开始以`}`结尾
func getNameAndRegstr(originStr string) (name, reg string) {
	pathVarStr := strings.TrimPrefix(originStr, "{")
	pathVarStr = strings.TrimSuffix(pathVarStr, "}")
	paramNameAndReg := utils.Split(pathVarStr, ":")
	name = paramNameAndReg[0]
	if len(paramNameAndReg) > 1 {
		reg = paramNameAndReg[1]
	}
	return
}

// getClass 根据参数获取参数类别
func getClass(s string) int {
	if s == suffixMatchSymbol {
		return matchAll
	}
	matchs, others := getMatchsAndOtehrs(s)
	matchsLen := len(matchs)
	othersLen := len(others)
	switch matchsLen {
	case 0:
		return normal
	case 1:
		if othersLen == 0 && !strings.ContainsRune(matchs[0], ':') {
			return pathVar
		}
		return pathReg

	default:
		if othersLen < matchsLen {
			panic("---")
		}
		return pathReg
	}
}

// addValueToPathParam 将参数值添加到map中
func addValueToPathParam(paramMap map[string]string, k, v string) {
	if paramMap != nil {
		paramMap[k] = v
	}
}

func getMatchsAndOtehrs(origin string) (matchs, others []string) {
	return utils.GetGroupMatch(origin, '{', '}')
}
