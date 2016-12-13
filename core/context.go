package core

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

/**
 *
 * @author doob  fudali113@gmail.com
 * @date 2016/12/7
 */

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	Params   map[string]string
}

func (this *Context) SetHttpStatus(num int) {
	this.response.WriteHeader(num)
}

func (this *Context) SetHeader(name, value string) {
	this.response.Header().Add(name, value)
}

func (this *Context) Param(name string) string {
	return this.Params[name]
}

func (this *Context) ParamInt(name string) int {
	strValue := this.request.Form.Get(name)
	value, _ := strconv.Atoi(strValue)
	return value
}

func (this *Context) BodyJson() string {
	body, err := ioutil.ReadAll(this.request.Body)
	if err != nil {
		logger.Error("get body str：%s", err.Error())
		return ""
	}
	return string(body)
}

func (this *Context) Write(bytes []byte) {
	this.response.Write(bytes)
}
