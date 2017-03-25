package errors

import (
	"fmt"
	"net/http"
)

// Matcher 是否匹配错误
type Matcher interface {
	Match(err interface{}) bool
}

// Dealer 匹配后的处理
type Dealer interface {
	Deal(err interface{}, w http.ResponseWriter, r *http.Request)
}

// ErrDealer 错误处理，包含匹配和处理
type ErrDealer interface {
	Matcher
	Dealer
}

// DoobError doob 的错误
type DoobError struct {
	Err       error
	Desc      string
	HttpStaus int
}

// Error 实现error
func (de DoobError) Error() string {
	return fmt.Sprintf(`request error : should return http status code is <%d> ,
		 description is <%s> , error is <%v>`, de.HttpStaus, de.Desc, de.Err)
}

// RequestError 请求错误
type RequestError struct {
	DoobError
	Url string
}

// Error 实现error
func (re RequestError) Error() string {
	return fmt.Sprintf("url is <%s> , error is %v", re.Url, re.DoobError)
}

// MehtodNotMatchError 匹配错误
type MehtodNotMatchError struct {
	RequestError
	Method string
}

// Error 实现error
func (mnm MehtodNotMatchError) Error() string {
	return fmt.Sprintf("url is <%s> ,method is <%s> , error is %v", mnm.Url, mnm.Method, mnm.DoobError)
}
