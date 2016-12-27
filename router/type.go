package router

import (
	"log"
	"regexp"
	"strings"

	"github.com/fudali113/doob/register"
)

// 返回值类型
type MatchResult struct {
	ParamMap map[string]string
	Rest     RestHandler
}

type RestHandler interface {
	// 获取该实列包含哪些method
	GetMethods() []string
	// 想该实列注入一个method的处理方法
	PutMethod(method string, handler register.RegisterHandlerType)
	// 根据方法名获取处理方法
	GetHandler(method string) register.RegisterHandlerType
	// 获取注册方法
	GetSigninHandler() register.RegisterHandlerType
	// 整合另一个RestHandler
	Joint(RestHandler)
}

// 实现 register package RegisterHandlerType 借口
type RegisterHandler struct {
	Handler      interface{}
	RegisterType *register.RegisterType
}

func (this *RegisterHandler) GetHandler() interface{} {
	return this.Handler
}
func (this *RegisterHandler) GetRegisterType() *register.RegisterType {
	return this.RegisterType
}

// RestHandler的简单实现
type SimpleRestHandler struct {
	methodHandlerMap map[string]register.RegisterHandlerType
}

// 获取一个 SimpleRestHandler 实例
func GetSimpleRestHandler(method string, sh interface{}) *SimpleRestHandler {
	registerHandler := &RegisterHandler{
		Handler:      sh,
		RegisterType: register.GetFuncRegisterType(sh),
	}
	return &SimpleRestHandler{
		methodHandlerMap: map[string]register.RegisterHandlerType{method: registerHandler},
	}
}

func (this *SimpleRestHandler) GetMethods() []string {
	res := make([]string, 0)
	for k, _ := range this.methodHandlerMap {
		res = append(res, k)
	}
	return res
}

func (this *SimpleRestHandler) PutMethod(method string, handler register.RegisterHandlerType) {
	this.methodHandlerMap[method] = handler
}
func (this *SimpleRestHandler) GetHandler(method string) register.RegisterHandlerType {
	res, ok := this.methodHandlerMap[method]
	if !ok {
		res = nil
	}
	return res
}
func (this *SimpleRestHandler) GetSigninHandler() register.RegisterHandlerType {
	for _, handler := range this.methodHandlerMap {
		if handler != nil {
			return handler
		}
	}
	log.Panic("注册函数不能为nil")
	return nil
}

func (this *SimpleRestHandler) Joint(restHandler RestHandler) {
	methods := restHandler.GetMethods()
	for _, method := range methods {
		this.PutMethod(method, restHandler.GetSigninHandler())
	}
}

// 实现Sort的接口
type nodes []*Node

func (this nodes) Len() int {
	return len(this)
}
func (this nodes) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this nodes) Less(i, j int) bool {
	a := this[i]
	b := this[j]
	return b.class > a.class
}

// 各类型储存接口
type nodeV interface {
	isMatch(urlPart string) (bool, bool)
	// if need pathvar
	// return in this method
	paramValue(urlPart string, url string) (bool, map[string]string)
	getOrigin() string
}

type nodeVNormal struct {
	origin string
}

func (this nodeVNormal) isMatch(urlPart string) (bool, bool) {
	return this.origin == urlPart, false
}
func (this nodeVNormal) paramValue(urlPart string, url string) (bool, map[string]string) {
	return false, nil
}
func (this nodeVNormal) getOrigin() string {
	return this.origin
}

type nodeVPathReg struct {
	origin    string
	paramName string
	paramReg  *regexp.Regexp
}

// check url part is match this node value
func (this nodeVPathReg) isMatch(urlPart string) (bool, bool) {
	findStr := this.paramReg.FindString(urlPart)
	log.Print(findStr, "====", urlPart)
	return findStr == urlPart, false
}
func (this nodeVPathReg) paramValue(urlPart string, url string) (bool, map[string]string) {
	return true, map[string]string{this.paramName: urlPart}
}
func (this nodeVPathReg) getOrigin() string {
	return this.origin
}

type nodeVPathVar struct {
	origin    string
	paramName string
}

func (this nodeVPathVar) isMatch(urlPart string) (bool, bool) {
	return true, false
}
func (this nodeVPathVar) paramValue(urlPart string, url string) (bool, map[string]string) {
	return true, map[string]string{this.paramName: urlPart}
}
func (this nodeVPathVar) getOrigin() string {
	return this.origin
}

type nodeVMatchAll struct {
	origin string
	prefix string
}

func (this nodeVMatchAll) isMatch(urlPart string) (bool, bool) {
	return strings.HasPrefix(urlPart, this.prefix), true
}
func (this nodeVMatchAll) paramValue(urlPart string, url string) (bool, map[string]string) {
	paramValue := strings.TrimPrefix(urlPart, this.prefix) + url
	return true, map[string]string{"**": paramValue}
}
func (this nodeVMatchAll) getOrigin() string {
	return this.origin
}
