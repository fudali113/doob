package errors

import "net/http"

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
