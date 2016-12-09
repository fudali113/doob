package router

import "github.com/fudali113/doob/core/register"

/**
 * 返回值类型
 */
type MatchResult struct {
	ParamMap     map[string]string
	Rest         RestHandler
	RegisterType *register.RegisterType
	ParamNames   []string
}
