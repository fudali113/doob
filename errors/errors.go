package errors

import (
	"net/http"
	"runtime/debug"

	. "github.com/fudali113/doob/http/const"
)

var (
	errDealers = []ErrDealer{}
)

// AddErrDealer 添加一个错误处理器
func AddErrDealer(errDealer ...ErrDealer) {
	errDealers = append(errDealers, errDealer...)
}

// CheckErr check panic err , match err and deal
func CheckErr(err interface{}, w http.ResponseWriter, r *http.Request, isDev bool) {

	defer func() {
		if e := recover(); e != nil {
			defaultErrDeal(e, w, isDev)
		}
	}()

	// traverse errDealers match true dealer and deal
	for _, errDealer := range errDealers {
		if errDealer.Match(err) {
			errDealer.Deal(err, w, r)
			return
		}
	}

	defaultErrDeal(err, w, isDev)

}

func defaultErrDeal(err interface{}, w http.ResponseWriter, isDev bool) {
	// default err dealer
	w.WriteHeader(INTERNAL_SERVER_ERROR)
	if isDev {
		WriteErrInfo(err, debug.Stack(), w)
	}
}
