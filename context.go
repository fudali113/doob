package doob

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	. "github.com/fudali113/doob/http_const"
	"io"
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
	return this.request.Form.Get(name)
}

func (this *Context) PostParam(name string) string {
	return this.request.PostForm.Get(name)
}

func (this *Context) PathParam(name string) string {
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

// Forward one request
func (this *Context) Forward(forwardUrl string, host ...string) {
	if len(host) == 0 {
		this.request.URL.Path = forwardUrl
		_doob.ServeHTTP(this.response, this.request)
		return
	}
	address := host[0] + forwardUrl
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
	this.response.Header().Del(CONTENT_TYPE)
	header := res.Header
	for k, v := range header {
		for _, v1 := range v {
			this.SetHeader(k, v1)
		}
	}
	body := make([]byte, 0)
	for {
		buf := make([]byte, redirectDefaultBodytLen)
		n, err := res.Body.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		body = append(body, buf[:n]...)
	}
	this.WriteBytes(body)
	this.response.WriteHeader(res.StatusCode)
}

// Redirect one request
func (this *Context) Redirect(redirectUrl string, host ...string) {
	address := func(host []string) string {
		if len(host) > 0 {
			return host[0]
		}
		return ""
	}(host) + redirectUrl
	this.SetHeader(LOCATION, address)
	if isDev {
		this.SetHeader(CACHE_CONTROL, NO_CACHE)
	}
	this.SetHttpStatus(MOVED_PERMANENTLY)
}
