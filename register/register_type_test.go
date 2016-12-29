package register

import (
	"github.com/fudali113/doob"
	"log"
	"testing"
)

func testFunc(string, string) string {
	return ""
}

func Test_GetFuncRegisterType(t *testing.T) {
	registerFunc := GetFuncRegisterType(testFunc)
	if registerFunc.ParamType.CiLen != 2 || registerFunc.ReturnType.Type != FILE {
		log.Print(registerFunc.ParamType.Type)
		log.Print(registerFunc.ParamType.CiLen)
		t.Error("GetFuncRegisterType have bug")
	}
}

func Test_getReturnType(t *testing.T) {
	returns := []string{"string"}
	returns1 := []string{"map[string]string"}
	returns2 := []string{"string", "map[string]string"}

	if getReturnType(returns).Type != FILE {
		t.Error("args")
	}
	if getReturnType(returns1).Type != JSON {
		t.Error("args")
	}
	if getReturnType(returns2).Type != RETURN_TYPE {
		t.Error("args")
	}
}

func Test_getParamType(t *testing.T) {
	params := []string{"string"}
	params1 := []string{"string", "string"}
	params2 := []string{"string", "string", "*doob.Context"}

	if getParamType(params).Type != CI_PATHVARIABLE || getParamType(params).CiLen != 1 {
		t.Error("args")
	}
	if getParamType(params1).Type != CI_PATHVARIABLE || getParamType(params1).CiLen != 2 {
		t.Error("args1")
	}
	if getParamType(params2).Type != CI_PATHVARIABLE_CTX || getParamType(params2).CiLen != 2 {
		t.Error("args2")
	}
}

func Test_ReturnType(t *testing.T) {
	test := func(*doob.Context) {}
	registerFunc := GetFuncRegisterType(test)
	if registerFunc.ReturnType.Type != RETURN_NONE {
		t.Error("dddddd")
	}
}
