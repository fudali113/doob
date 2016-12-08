package reflect

import (
	"testing"

	"github.com/fudali113/doob/core"
)

func test(string, string, string, *core.Context) string {
	return ""
}

func Test_GetFuncParams(t *testing.T) {
	params, returns := GetFuncParams(test)
	logger.Debug("%v", params)
	logger.Debug("%v", returns)
	if len(params) != 4 {
		t.Error("GetFuncParams have bug")
	}
	if len(returns) != 1 {
		t.Error("GetFuncParams have bug")
	}
}
