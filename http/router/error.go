package router

// NotMatch 未匹配错误
type NotMatch struct {
	errorStr string
}

func (nm NotMatch) Error() string {
	return nm.errorStr
}
