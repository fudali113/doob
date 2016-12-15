package return_deal

import (
	"strings"

	"github.com/fudali113/doob/utils"
)

// 判断是否前匹配
func matchPrefix(str string, prefixs ...string) bool {
	res := false
	str = strings.ToLower(str)
	for _, prefix := range prefixs {
		if strings.HasPrefix(str, prefix) {
			res = true
			break
		}
	}
	return res
}

// 获取 typeStr 中的路径
func getPath(typeStr string) string {
	typeAndPath := utils.Split(typeStr, ":")
	if len(typeAndPath) == 2 {
		path := typeAndPath[1]
		suffix := "." + typeAndPath[0]
		// 如果文件没有跟上后缀名
		// 加上后缀名
		if !strings.HasSuffix(path, suffix) {
			path = path + suffix
		}
		return path
	}
	return ""
}
