package doob

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	. "github.com/fudali113/doob/http_const"
)

// 对 ResponseWriter 和 request 封装的上下文
type Context struct {
	request    *http.Request
	response   http.ResponseWriter
	pathParams map[string]string
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
	return this.pathParams[name]
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
		log.Print("get body str : ", err.Error())
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
	this.SetHeader(CONTENT_TYPE, APP_JSON)
}

// redirect
// no test
func (this *Context) Redirect(redirectUrl string, addresses ...string) {
	if len(addresses) == 0 {
		this.request.URL.Path = redirectUrl
		_doob.ServeHTTP(this.response, this.request)
		return
	}
	address := addresses[0] + redirectUrl
	client := &http.Client{}
	request, err := http.NewRequest(this.request.Method, address, this.request.Body)
	if err != nil {
		log.Print("Redirect is error , error is ", err)
		this.SetHttpStatus(INTERNAL_SERVER_ERROR)
		return
	}
	res, err := client.Do(request)
	if err != nil {
		log.Print("Redirect is error , error is ", err)
		this.SetHttpStatus(INTERNAL_SERVER_ERROR)
		return
	}
	this.response.WriteHeader(res.StatusCode)
	header := res.Header
	for k, v := range header {
		for _, v1 := range v {
			this.response.Header().Add(k, v1)
		}
	}
	reader := io.TeeReader(res.Body, this.response)
	body := make([]byte, redirectDefaulBodytLen)
	reader.Read(body)
}
