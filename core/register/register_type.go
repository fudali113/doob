package register

import (
	"log"

	"github.com/fudali113/doob/utils/reflect"
)

/**
 * 函数类别
 */
const (
	PARAM_NONE = iota
	CTX
	ORIGIN
	CI_PATHVARIABLE
	CI_PATHVARIABLE_CTX
	CI_PATHVARIABLE_ORIGIN
)

const (
	RETURN_NONE = iota
	JSON        = iota
	FILE
	MODAL_HTML
	HTML_MODAL
)

type RegisterType struct {
	Handler    interface{}
	ParamType  *ParamType
	ReturnType *ReturnType
}

type ParamType struct {
	Type  int
	CiLen int
}

type ReturnType struct {
	Type int
}

func GetFuncRegisterType(function interface{}) *RegisterType {
	paramType, returnType := GetFuncParam3ReturnType(function)
	return &RegisterType{
		Handler:    function,
		ParamType:  paramType,
		ReturnType: returnType,
	}
}

func GetFuncParam3ReturnType(function interface{}) (*ParamType, *ReturnType) {
	params, returns := reflect.GetFuncParams(function)
	return getParamType(params), getReturnType(returns)
}

/**
 * 获取参数类型
 */
func getParamType(params []string) *ParamType {
	stringTypeLen := 0
	hasCTX := 0
	hasOringin := 0
	for _, param := range params {
		log.Print(param)
		switch param {
		case "string":
			stringTypeLen++
			if hasOringin > 0 || hasCTX > 0 {
				log.Panic("自动注入url参数必须放在参数最前面")
			}
		case "*core.Context":
			hasCTX++
		case "*http.Request":
			hasOringin++
		case "http.ResponseWriter":
			hasOringin++
		default:
		}
	}
	if stringTypeLen > 0 {
		if hasCTX > 0 {
			return &ParamType{
				Type:  CI_PATHVARIABLE_CTX,
				CiLen: stringTypeLen,
			}
		}
		if hasOringin > 0 {
			return &ParamType{
				Type:  CI_PATHVARIABLE_ORIGIN,
				CiLen: stringTypeLen,
			}
		}
		return &ParamType{
			Type:  CI_PATHVARIABLE,
			CiLen: stringTypeLen,
		}
	}
	if hasCTX > 0 {
		return &ParamType{
			Type: CTX,
		}
	}
	if hasOringin > 0 {
		return &ParamType{
			Type: ORIGIN,
		}
	}
	return &ParamType{
		Type: PARAM_NONE,
	}
}

/**
 * 获取返回值类型
 */
func getReturnType(returns []string) *ReturnType {
	Type := func(returns []string) int {
		switch len(returns) {
		case 1:
			switch returns[0] {
			case "string":
				return FILE
			default:
				return JSON

			}
		case 2:
			switch returns[0] {
			case "string":
				return HTML_MODAL
			default:
				return MODAL_HTML
			}
		default:
			return RETURN_NONE
		}
	}(returns)

	return &ReturnType{Type: Type}
}
