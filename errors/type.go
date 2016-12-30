package errors

import (
	"fmt"
	"net/http"
)

type Error int

func (err Error) Error() string {
	return fmt.Sprintf("weight value is %d", int(err))
}

func (err Error) GetWV() int {
	return int(err)
}

type Matcher interface {
	Match(err interface{}) bool
}

type Dealer interface {
	Deal(err interface{}, w http.ResponseWriter, r *http.Request)
}

type ErrDealer interface {
	Matcher
	Dealer
}
