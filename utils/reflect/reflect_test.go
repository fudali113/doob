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

func test1(name string) string {
	return name
}

func Test_Invoke(t *testing.T) {
	log.Print("oooooooo___________", Invoke(test1, "name"))
}
