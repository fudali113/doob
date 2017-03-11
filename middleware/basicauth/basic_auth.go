package basic_auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/fudali113/doob/config"
	"github.com/fudali113/doob/middleware"

	. "github.com/fudali113/doob/http/const"
)

const (
	BASIC_AUTH_PREFIX = "Basic "
)

// impl basic auth
func BasicAuth(res http.ResponseWriter, req *http.Request) (ispass bool) {
	authStr := req.Header.Get(BASIC_AUTH)
	if config.OpenBasicAuth {
		if len(authStr) > 0 && strings.HasPrefix(authStr, BASIC_AUTH_PREFIX) {
			authBase64Str := strings.TrimPrefix(authStr, BASIC_AUTH_PREFIX)
			username, passwd, err := getUsernameAndPasswd(authBase64Str)
			if err != nil {
				ispass = false
				goto deal
			}
			if config.BasicAuthUserInConfig[username] == passwd {
				ispass = true
			} else {
				ispass = false
			}
		} else {
			ispass = false
		}
	deal:
		if !ispass {
			res.Header().Add(WWW_AUTH, `Basic realm="`+config.BasicAuthReminder+`"`)
			res.WriteHeader(UNAUTHORIZED)
		}
		return
	}
	ispass = true
	return
}

// from basic auth base64 string
// get username and passwd
func getUsernameAndPasswd(authBase64Str string) (username, passwd string, _ error) {
	userAndPasswdBytes, err := base64.StdEncoding.DecodeString(authBase64Str)
	if err != nil {
		return "", "", err
	}
	userAndPasswdStr := string(userAndPasswdBytes)
	userAndPasswd := strings.Split(userAndPasswdStr, ":")
	if len(userAndPasswd) != 2 {
		return "", "", fmt.Errorf("basic auth string format error")
	}
	return userAndPasswd[0], userAndPasswd[1], nil
}

func init() {
	basicAuhtFilter := middleware.BeforeFilterFunc(BasicAuth)
	middleware.AddBFilter(&basicAuhtFilter)
}
