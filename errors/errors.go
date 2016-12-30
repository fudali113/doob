package errors

import (
	"net/http"
	"runtime/debug"

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
		if e := recover(); e != nil {
			// default err dealer
			w.WriteHeader(INTERNAL_SERVER_ERROR)
			if isDev {
				WriteErrInfo(e, debug.Stack(), w)
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
