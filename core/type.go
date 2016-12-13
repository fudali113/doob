package core

// 自动注入参数function兼容
// 不符合这些条件的func讲使用反射执行方法
//
type ReturnStr func() string
type ReturnObject func() interface{}
type ReturnType func() (string, interface{})

type CTXReturn func(*Context)
type CTXReturnStr func(*Context) string

type ctxReturnObject interface {
	ServerHTTP(ctx *Context) interface{}
}

type CTXReturnObject func(*Context) interface{}

func (c CTXReturnObject) ServerHTTP(ctx *Context) interface{} {
	return c(ctx)
}

type CTXReturnType func(*Context) (string, interface{})

// =====================================================
// 下面想不用反射完成自动注入参数，但是貌似工作量太大了-_-
// =====================================================

type DI1ReturnStr func(string) string
type DI2ReturnStr func(string, string) string
type DI3ReturnStr func(string, string, string) string

type DI1ReturnObject func(string) interface{}
type DI2ReturnObject func(string, string) interface{}
type DI3ReturnObject func(string, string, string) interface{}

type DI1CTXReturnStr func(string, *Context) string
type DI2CTXReturnStr func(string, string, *Context) string
type DI3CTXReturnStr func(string, string, string, *Context) string

type DI1CTXReturnObject func(string, *Context) interface{}
type DI2CTXReturnObject func(string, string, *Context) interface{}
type DI3CTXReturnObject func(string, string, string, *Context) interface{}
