package register

import (
	"log"
	"net/http"
	"testing"
)

func testFunc(string, string, http.ResponseWriter, *http.Request) string {
	return ""
}

func Test_GetFuncRegisterType(t *testing.T) {
	registerFunc := GetFuncRegisterType(testFunc)
	if registerFunc.ParamType.CiLen != 2 ||
		registerFunc.ParamType.Type != CI_PATHVARIABLE_ORIGIN ||
		registerFunc.ReturnType.Type != FILE {
		log.Print(registerFunc.ParamType.Type)
		log.Print(registerFunc.ParamType.CiLen)
		t.Error("GetFuncRegisterType have bug")
	}
}
