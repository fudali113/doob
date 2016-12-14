package router

import (
	"log"

	"github.com/fudali113/doob/core/register"
)

type Router interface {
	// 添加一个处理器
	Add(string, RestHandler)
	// 根据url获取一个最佳处理器
	Get(string) *MatchResult
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
