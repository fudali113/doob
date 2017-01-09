package session

import (
	"net/http"
	"reflect"

	"github.com/fudali113/doob/config"
	"github.com/fudali113/doob/middleware"
	"github.com/fudali113/doob/utils"

	. "github.com/fudali113/doob/http/const"
)

const (
	cookieName = "doob_id"
)

var (
	session = &SessionMW{Repo: map[string]Session{}}

	CreateSeesionCookieValueFunc = func() string {
		return utils.GetMd5String(utils.GetRandomStr(), config.SessionCreateSecretKey)
	}
)

func GetSession(req *http.Request) (Session, error) {
	cookie, err := req.Cookie(cookieName)
	var cookieV string
	if err != nil {
		return nil, err
	} else {
		cookieV = cookie.Value
	}
	thisSession := session.Repo[cookieV]
	if thisSession == nil {
		oo := sessionMemryRepo(map[string]interface{}{})
		thisSession = &oo
		session.Repo[cookieV] = thisSession
	}
	return thisSession, nil
}

type SessionMW struct {
	Repo map[string]Session
}

func (this SessionMW) DoBeforeFilter(w http.ResponseWriter, req *http.Request) (isPass bool) {

	isPass = true

	cookie, err := req.Cookie(cookieName)
	var cookieV string
	if err != nil {
		cookieV = ""
	} else {
		cookieV = cookie.Value
	}
	if cookieV == "" || this.Repo[cookieV] == nil {
		cookieV = CreateSeesionCookieValueFunc()
		thisSession := sessionMemryRepo(map[string]interface{}{})
		this.Repo[cookieV] = &thisSession
		w.Header().Add(SET_COOKIE, cookieName+"="+cookieV)
		return
	}
	return
}

func (this SessionMW) DoLaterFilter(res http.ResponseWriter, req *http.Request) {

}

type Session interface {
	Set(string, interface{})
	Get(string) interface{}
	GetByPointer(string, interface{})
}

type sessionMemryRepo map[string]interface{}

func (this sessionMemryRepo) Set(k string, v interface{}) {
	this[k] = v
}

func (this sessionMemryRepo) Get(k string) interface{} {
	return this[k]
}

func (this sessionMemryRepo) GetByPointer(k string, vPointer interface{}) {
	v := this[k]
	vType := reflect.TypeOf(&v)
	vPointerType := reflect.TypeOf(vPointer)
	if vType == vPointerType {
		vPointer = &v
	} else {
		panic("oo")
	}
}

func init() {
	middleware.AddMiddlerware(session)
}
