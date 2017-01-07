package http

import (
	"github.com/fudali113/doob/http/router"
	"github.com/fudali113/doob/middleware"
)

var (
	beforeFilters = make([]middlerware.BeforeFilter, 0)
	laterFilters  = make([]middlerware.LaterFilter, 0)
	root          = router.GetRoot()

	doob = &Doob{
		bFilters: beforeFilters,
		lFilters: laterFilters,
		Root:     root,
	}
)

func GetDoob() *Doob {
	return doob
}
