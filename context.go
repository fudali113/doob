package doob

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/fudali113/doob/config"
	"github.com/fudali113/doob/errors"
	"github.com/fudali113/doob/middleware/session"

	. "github.com/fudali113/doob/http/const"
)

// 对 ResponseWriter 和 request 封装的上下文
type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	PathParams map[string]string
}

// 设置 http 返回状态码
func (this *Context) SetHttpStatus(num int) {
	this.Response.WriteHeader(num)
}

// 添加返回header
func (this *Context) AddHeader(name, value string) {
	this.Response.Header().Add(name, value)
}

// 获取 参数名为 name 的参数值
func (this *Context) Param(name string) string {
	return this.Request.Form.Get(name)
}

func (this *Context) PostParam(name string) string {
	return this.Request.PostForm.Get(name)
}

func (this *Context) PathParam(name string) string {
	return this.PathParams[name]
}

func (this *Context) Seesion() session.Session {
	s, err := session.GetSession(this.Request)
	if err != nil {
		return nil
	}
	return s
}

// 获取参数名为 name 的参数值并转化为int类型
// 当转化失败时返回 0
func (this *Context) ParamInt(name string) int {
	strValue := this.Request.Form.Get(name)
	value, _ := strconv.Atoi(strValue)
	return value
}

func (this *Context) Ip() string {
	return this.Request.RemoteAddr
}

func (this *Context) URI() string {
	return this.Request.RequestURI
}

// 获取请求的 body 字符串
func (this *Context) BodyJson() string {
	body, err := ioutil.ReadAll(this.Request.Body)
	if err != nil {
		log.Print("get body str : ", err.Error())
		return ""
	}
	return string(body)
}

// 往 responseWriter 中写入内容
func (this *Context) WriteBytes(bytes []byte) {
	this.Response.Write(bytes)
}

// 接受一个实体并转化为 json bytes
// 往 responseWriter 中写入 改 json bytes
// 并添加 header Application/json
func (this *Context) WriteJson(jsonStruct interface{}) {
	json, err := json.Marshal(jsonStruct)
	if err != nil {
		return
	}
	this.Response.Write(json)
	this.AddHeader(CONTENT_TYPE, APP_JSON)
}

// Forward one request
//
// @panic DoobError
func (this *Context) Forward(forwardUrl string, host ...string) {
	if len(host) == 0 {
		this.Request.URL.Path = forwardUrl
		doob.ServeHTTP(this.Response, this.Request)
		return
	}
	address := host[0] + forwardUrl
	client := &http.Client{}
	request, err := http.NewRequest(this.Request.Method, address, this.Request.Body)
	if err != nil {
		panic(errors.DoobError{
			Err:       err,
			Desc:      "forward : create request is error",
			HttpStaus: INTERNAL_SERVER_ERROR,
		})
	}
	res, err := client.Do(request)
	if err != nil {
		panic(errors.DoobError{
			Err:       err,
			Desc:      "forward : do request is error",
			HttpStaus: INTERNAL_SERVER_ERROR,
		})
	}
	this.Response.Header().Del(CONTENT_TYPE)
	header := res.Header
	for k, v := range header {
		for _, v1 := range v {
			this.AddHeader(k, v1)
		}
	}
	body := make([]byte, 0)
	for {
		buf := make([]byte, config.RedirectDefaultBodyLen)
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
	this.Response.WriteHeader(res.StatusCode)
}

// Redirect one request
func (this *Context) Redirect(redirectUrl string, host ...string) {
	address := func(host []string) string {
		if len(host) > 0 {
			return host[0]
		}
		return ""
	}(host) + redirectUrl
	this.AddHeader(LOCATION, address)
	if config.IsDev {
		this.AddHeader(CACHE_CONTROL, NO_CACHE)
	}
	this.SetHttpStatus(MOVED_PERMANENTLY)
}
