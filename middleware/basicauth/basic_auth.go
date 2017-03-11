package basic_auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/fudali113/doob/config"
	"github.com/fudali113/doob/middleware"

	. "github.com/fudali113/doob/http/const"
)

func BasicAuth(res http.ResponseWriter, req *http.Request) (ispass bool) {
	authStr := req.Header.Get(BASIC_AUTH)
	if config.OpenBasicAuth {
		if len(authStr) > 0 && strings.HasPrefix(authStr, "Basic ") {
			authBase64Str := strings.TrimPrefix(authStr, "Basic ")
			userAndPasswdBytes, err := base64.StdEncoding.DecodeString(authBase64Str)
			if err != nil {
				ispass = false
				goto deal
			}
			userAndPasswdStr := string(userAndPasswdBytes)
			userAndPasswd := strings.Split(userAndPasswdStr, ":")
			if len(userAndPasswd) != 2 {
				ispass = false
				goto deal
			}
			user := userAndPasswd[0]
			passwd := userAndPasswd[1]
			if config.BasicAuthUserInConfig[user] == passwd {
				ispass = true
			} else {
				ispass = false
			}
		} else {
			ispass = false
		}
	deal:
		if !ispass {
			res.WriteHeader(UNAUTHORIZED)
		}
		return
	} else {
		ispass = true
	}
	return
}

func init() {
	middleware.AddBFilter(middleware.BeforeFilterFunc(BasicAuth))
}
