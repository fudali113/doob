package reflect

import (
	"reflect"

	"regexp"
	"strings"

	"github.com/fudali113/doob/log"
	"github.com/fudali113/doob/utils"
)

var (
	logger      = log.GetLog("reflect")
	paramReg, _ = regexp.Compile("\\([\\S|\\s]+?\\)")
)

func GetFuncParams(function interface{}) (params []string, returns []string) {
	funcType := reflect.TypeOf(function)
	funcStr := funcType.String()

	params3returns := paramReg.FindAllString(funcStr, -1)
	if len(params3returns) == 2 {
		params = getSlice(params3returns[0])
		returns = getSlice(params3returns[1])
	}

	switch len(params3returns) {
	case 1:
		params = getSlice(params3returns[0])
		funcStr = strings.Replace(funcStr, "func", "", -1)
		funcStr = paramReg.ReplaceAllString(funcStr, "")
		returns = append(returns, funcStr)
	case 2:
		params = getSlice(params3returns[0])
		returns = getSlice(params3returns[1])
	default:

	}
	return nil, nil
}

func getSlice(arg string) []string {
	res := utils.Split(arg[1:len(arg)-1], ",")
	return res
}
