package errors

import (
	"net/http"
	"runtime/debug"

	"errors"
	. "github.com/fudali113/doob/http_const"
)

var (
	errDealers = []ErrDealer{}
)

func AddErrDealer(errDealer ...ErrDealer) {
	errDealers = append(errDealers, errDealer...)
}

// check panic err , match err and deal
func CheckErr(err interface{}, w http.ResponseWriter, r *http.Request, isDev bool) {

	defer func() {
		if err := recover(); err != nil {
			// default err dealer
			w.WriteHeader(INTERNAL_SERVER_ERROR)
			if isDev {
				WriteErrInfo(err, debug.Stack(), w)
			}
		}
	}()

	// traverse errDealers match true dealer and deal
	for _, errDealer := range errDealers {
		if errDealer.Match(err) {
			errDealer.Deal(err, w, r)
			return
		}
	}

}
