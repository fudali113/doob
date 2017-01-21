package reflect

import (
	"log"
	"net/http"
	"testing"
)

func test(string, string, string, http.ResponseWriter) string {
	return ""
}

func Test_GetFuncParams(t *testing.T) {
	params, returns := GetFuncParams(test)
	if len(params) != 4 {
		t.Error("GetFuncParams have bug")
	}
	log.Print(returns)
	if len(returns) != 1 {
		t.Error("GetFuncParams have bug")
	}
}

func Test_ContrastType(t *testing.T) {
	type test struct {
		oo string
		aa string
	}

	a := test{"a", "b"}
	b := test{"b", "a"}
	if !ContrastType(a, b) {
		t.Error("ContrastType func has bug")
	}

	var a1 interface{} = test{"a", "b"}
	b1 := &test{"b", "a"}
	if !ContrastType(&a1, b1) {
		t.Error("ContrastType func has bug")
	}
}

func test1(name string) string {
	return name
}

func Test_Invoke(t *testing.T) {
	log.Print("oooooooo___________", Invoke(test1, "name"))
}
