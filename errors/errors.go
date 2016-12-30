package errors

import (
	"fmt"
	"runtime/debug"
	"net/http"

	. "github.com/fudali113/doob/http_const"
)

// check panic err , match err and deal
func CheckErr(err interface{}, w http.ResponseWriter, isDev bool)  {
	switch err.(type) {
	default:
		w.WriteHeader(INTERNAL_SERVER_ERROR)
		if isDev {
			WriteErrInfo(err,debug.Stack(),w)
		}
	}
}

type Error int

func (err Error) Error() string {
	return fmt.Sprintf("weight value is %d", int(err))
}

func (err Error) GetWV() int {
	return int(err)
}

func GetMwthodMatchError(should, fact, url string, errors ...error) *MethodMacthError {
	return &MethodMacthError{
		shouldMethod: should,
		factMethod:   fact,
		matchError: &URLMacthError{
			url: url,
			matchError: &MatchError{
				message: "",
			},
		},
	}
}

func GetURLMatchError(url, message string, errors ...error) *URLMacthError {
	return &URLMacthError{
		url: url,
		matchError: &MatchError{
			message: message,
		},
	}
}
