package core

import (
	"log"
	"regexp"
	"strings"

	"github.com/fudali113/doob/core/errors"
)

func (this *urlInfo) addUrlPara(v urlMacthPara) {
	this.urlParas = append(this.urlParas, v)
}

func (this *urlInfo) len() int {
	return len(this.urlParas)
}

func convertHttpMethods2String(tms []HttpMethod) (string, error) {
	res := []string{}
	hasRightMethodStr := false
	for i := 0; i < len(tms); i++ {
		methodStr := strings.ToLower(string(tms[i]))
		match, _ := regexp.MatchString("get|post|put|delete|options|head", methodStr)
		if match {
			res = append(res, string(tms[i]))
			hasRightMethodStr = true
		} else {
			log.Printf("you have a http method is barbarism : %s ", methodStr)
		}
	}
	if !hasRightMethodStr {
		return "", &errors.MethodStrBarbarismError{}
	}
	return strings.Join(res, ","), nil
}
