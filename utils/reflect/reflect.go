package reflect

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/fudali113/doob/utils"
)

var (
	paramReg, _ = regexp.Compile("\\([\\S|\\s]*?\\)")
)

func Invoke(function interface{}, params ...interface{}) []interface{} {
	valueParams := []reflect.Value{}
	for _, param := range params {
		value := reflect.ValueOf(param)
		valueParams = append(valueParams, value)
	}
	funcValue := reflect.ValueOf(function)
	values := funcValue.Call(valueParams)
	returns := []interface{}{}
	for _, value := range values {
		returns = append(returns, value.Interface())
	}
	return returns
}

// 获取方法的参数和返回值类型数组
func GetFuncParams(function interface{}) (params []string, returns []string) {
	funcType := reflect.TypeOf(function)
	funcStr := funcType.String()
	funcStr = strings.Replace(funcStr, " ", "", -1)

	params3returns := paramReg.FindAllString(funcStr, -1)
	switch len(params3returns) {
	case 1:
		params = getSlice(params3returns[0])
		funcParamReg, _ := regexp.Compile("func\\([\\S|\\s]*\\)")
		funcStr = funcParamReg.ReplaceAllString(funcStr, "")
		returns = append(returns, funcStr)
	case 2:
		params = getSlice(params3returns[0])
		returns = getSlice(params3returns[1])
	default:
	}
	return
}

// 比较两个类型是否相同
func ContrastType(a, b interface{}) (eq bool) {
	defer func() {
		if err := recover(); err != nil {
			eq = false
		}
	}()
	aType := reflect.TypeOf(a).Name()
	bType := reflect.TypeOf(b).Name()
	eq = aType == bType
	return
}

func getSlice(arg string) []string {
	if strings.Replace(arg, " ", "", -1) == "" {
		return []string{}
	}
	res := utils.Split(arg[1:len(arg)-1], ",")
	return res
}
