package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 对 ResponseWriter 和 request 封装的上下文
type Context struct {
	request  *http.Request
	response http.ResponseWriter
	Params   map[string]string
}

// 设置 http 返回状态码
func (this *Context) SetHttpStatus(num int) {
	this.response.WriteHeader(num)
}

// 添加返回header
func (this *Context) SetHeader(name, value string) {
	this.response.Header().Add(name, value)
}

// 获取 参数名为 name 的参数值
func (this *Context) Param(name string) string {
	return this.Params[name]
}

// 获取参数名为 name 的参数值并转化为int类型
// 当转化失败时返回 0
func (this *Context) ParamInt(name string) int {
	strValue := this.request.Form.Get(name)
	value, _ := strconv.Atoi(strValue)
	return value
}

// 获取请求的 body 字符串
func (this *Context) BodyJson() string {
	body, err := ioutil.ReadAll(this.request.Body)
	if err != nil {
		logger.Error("get body str：%s", err.Error())
		return ""
	}
	return string(body)
}

// 往 responseWriter 中写入内容
func (this *Context) WriteBytes(bytes []byte) {
	this.response.Write(bytes)
}

// 接受一个实体并转化为 json bytes
// 往 responseWriter 中写入 改 json bytes
// 并添加 header Application/json
func (this *Context) WriteJson(jsonStruct interface{}) {
	json, err := json.Marshal(jsonStruct)
	if err != nil {
		return
	}
	this.response.Write(json)
	this.SetHeader("context", "application/json")
}
