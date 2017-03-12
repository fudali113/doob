package router

import (
	"log"
	"regexp"
	"strings"

	"github.com/fudali113/doob/register"
)

// MatchResult 返回值类型
// 用于装在获取的返回值
type MatchResult struct {
	ParamMap map[string]string
	Rest     RestHandler
}

// RestHandler 同于存储不同方法的不同处理器
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

// RegisterHandler 实现 register package RegisterHandlerType 借口
type RegisterHandler struct {
	Handler      interface{}
	RegisterType *register.RegisterType
}

// GetHandler 获取handler
func (rh *RegisterHandler) GetHandler() interface{} {
	return rh.Handler
}

// GetRegisterType 获取注册类型
func (rh *RegisterHandler) GetRegisterType() *register.RegisterType {
	return rh.RegisterType
}

// SimpleRestHandler RestHandler的简单实现
type SimpleRestHandler struct {
	methodHandlerMap map[string]register.RegisterHandlerType
}

// GetSimpleRestHandler 获取一个 SimpleRestHandler 实例
func GetSimpleRestHandler(method string, sh interface{}) *SimpleRestHandler {
	registerHandler := &RegisterHandler{
		Handler:      sh,
		RegisterType: register.GetFuncRegisterType(sh),
	}
	return &SimpleRestHandler{
		methodHandlerMap: map[string]register.RegisterHandlerType{method: registerHandler},
	}
}

// GetMethods 获取有哪些http方法
func (srh *SimpleRestHandler) GetMethods() []string {
	res := make([]string, 0)
	for k := range srh.methodHandlerMap {
		res = append(res, k)
	}
	return res
}

// PutMethod add a method handler
func (srh *SimpleRestHandler) PutMethod(method string, handler register.RegisterHandlerType) {
	srh.methodHandlerMap[method] = handler
}

// GetHandler get a method handler with method
func (srh *SimpleRestHandler) GetHandler(method string) register.RegisterHandlerType {
	res, ok := srh.methodHandlerMap[method]
	if !ok {
		res = nil
	}
	return res
}

// GetSigninHandler 获取注册方法
func (srh *SimpleRestHandler) GetSigninHandler() register.RegisterHandlerType {
	for _, handler := range srh.methodHandlerMap {
		if handler != nil {
			return handler
		}
	}
	log.Panic("注册函数不能为nil")
	return nil
}

// Joint 合并另外一个RestHandler
func (srh *SimpleRestHandler) Joint(restHandler RestHandler) {
	methods := restHandler.GetMethods()
	for _, method := range methods {
		srh.PutMethod(method, restHandler.GetSigninHandler())
	}
}

// 实现Sort的接口
type nodes []*Node

func (n nodes) Len() int {
	return len(n)
}
func (n nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
func (n nodes) Less(i, j int) bool {
	a := n[i]
	b := n[j]
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

func (nvn nodeVNormal) isMatch(urlPart string) (bool, bool) {
	return nvn.origin == urlPart, false
}
func (nvn nodeVNormal) paramValue(urlPart string, url string) (bool, map[string]string) {
	return false, nil
}
func (nvn nodeVNormal) getOrigin() string {
	return nvn.origin
}

type nodeVPathReg struct {
	origin    string
	paramName string
	paramReg  *regexp.Regexp
}

// check url part is match nvpg node value
func (nvpg nodeVPathReg) isMatch(urlPart string) (bool, bool) {
	findStr := nvpg.paramReg.FindString(urlPart)
	return findStr == urlPart, false
}
func (nvpg nodeVPathReg) paramValue(urlPart string, url string) (bool, map[string]string) {
	return true, map[string]string{nvpg.paramName: urlPart}
}
func (nvpg nodeVPathReg) getOrigin() string {
	return nvpg.origin
}

type nodeVPathVar struct {
	origin    string
	paramName string
}

func (nvpv nodeVPathVar) isMatch(urlPart string) (bool, bool) {
	return true, false
}
func (nvpv nodeVPathVar) paramValue(urlPart string, url string) (bool, map[string]string) {
	return true, map[string]string{nvpv.paramName: urlPart}
}
func (nvpv nodeVPathVar) getOrigin() string {
	return nvpv.origin
}

type nodeVMatchAll struct {
	origin string
	prefix string
}

func (nvma nodeVMatchAll) isMatch(urlPart string) (bool, bool) {
	return strings.HasPrefix(urlPart, nvma.prefix), true
}
func (nvma nodeVMatchAll) paramValue(urlPart string, url string) (bool, map[string]string) {
	paramValue := strings.TrimPrefix(urlPart, nvma.prefix) + url
	return true, map[string]string{"**": paramValue}
}
func (nvma nodeVMatchAll) getOrigin() string {
	return nvma.origin
}
