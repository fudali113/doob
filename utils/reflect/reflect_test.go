package reflect

import (
	"net/http"
	"testing"
)

func test(string, string, string, http.ResponseWriter) string {
	return ""
}

func returnHtml() string {
	return "static/index.html"
}

func Test_GetFuncParams(t *testing.T) {
	params, returns := GetFuncParams(returnHtml)
	logger.Debug("%v", params)
	logger.Debug("%v", returns)
	if len(params) != 4 {
		t.Error("GetFuncParams have bug")
	}
	logger.Debug("------ %v", returns)
	if len(returns) != 1 {
		t.Error("GetFuncParams have bug")
	}
}

func test1(name string) string {
	return name
}

func Test_Invoke(t *testing.T) {
	logger.Debug("oooooooo___________%v", Invoke(test1, "name"))
}
