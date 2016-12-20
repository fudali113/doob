package tree_router

import (
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

// create a new node
func creatNode(url string, rt reserveType) (newNode *node, isOver bool) {
	prefix, other := splitUrl(url)
	isOver = false
	newNode = &node{
		class: getClass(prefix),
		value: prefix,
	}
	if strings.TrimSpace(other) == "" {
		newNode.handler = rt
		isOver = true
	} else {
		newNode.children = make([]*node, 0)
	}
	return
}

func isMatch(fact , origin string, class int) bool  {
	switch class {
	case normal:
		return fact == origin
	case pathReg:
		return
	}
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
