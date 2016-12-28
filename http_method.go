package doob

import "regexp"

type HttpMethod string

func checkHttpMethod(httpMethod HttpMethod) bool {
	return checkMethodStr(string(httpMethod))
}

func checkMethodStr(httpMethod string) bool {
	match, _ := regexp.MatchString("get|post|put|delete|options|head", httpMethod)
	return match
}
