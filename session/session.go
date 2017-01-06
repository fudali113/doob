package session

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"reflect"

	. "github.com/fudali113/doob/http_const"
	"github.com/fudali113/doob/middleware"
)

const (
	cookieName = "doob_id"
)

var (
	session = &SessionMW{Repo: map[string]Session{}}
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

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
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
		cookieV = func() string {
			b := make([]byte, 48)

			if _, err := io.ReadFull(rand.Reader, b); err != nil {
				return ""
			}
			return GetMd5String(base64.URLEncoding.EncodeToString(b))
		}()
		this.Repo[cookieV] = sessionMemryRepo(map[string]interface{}{})
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
	middlerware.AddMiddlerware(session)
}
