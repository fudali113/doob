package router

// 返回值类型
type MatchResult struct {
	ParamMap   map[string]string
	Rest       RestHandler
	ParamNames []string
}
