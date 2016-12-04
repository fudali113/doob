package router

type Router interface {
	/**
	 * 添加一个处理器
	 */
	Add(string, RestHandler)
	/**
	 * 根据url获取一个最佳处理器
	 */
	Get(string) *MatchResult
}

type RestHandler interface {
	/**
	 * 获取该实列包含哪些method
	 */
	GetMethods() []string
	/**
	 * 想该实列注入一个method的处理方法
	 */
	PutMethod(method string, handler interface{})
	/**
	 * 根据方法名获取处理方法
	 */
	GetHandler(method string) interface{}
	/**
	 * 获取注册方法
	 */
	GetSigninHandler() interface{}
	/**
	 * 整合另一个RestHandler
	 */
	Joint(RestHandler)
}

/**
 * RestHandler的简单实现
 */
type SimpleRestHandler struct {
	methodHandlerMap map[string]interface{}
	signinHandler    interface{}
}

func GetSimpleRestHandler(mhm map[string]interface{}, sh interface{}) *SimpleRestHandler {
	return &SimpleRestHandler{
		methodHandlerMap: mhm,
		signinHandler:    sh,
	}
}

func (this *SimpleRestHandler) GetMethods() []string {
	res := make([]string, 0)
	for k, _ := range this.methodHandlerMap {
		res = append(res, k)
	}
	return res
}

func (this *SimpleRestHandler) PutMethod(method string, handler interface{}) {
	this.methodHandlerMap[method] = handler
}
func (this *SimpleRestHandler) GetHandler(method string) interface{} {
	res, ok := this.methodHandlerMap[method]
	if !ok {
		res = nil
	}
	return res
}
func (this *SimpleRestHandler) GetSigninHandler() interface{} {
	return this.signinHandler
}

func (this *SimpleRestHandler) Joint(restHandler RestHandler) {
	methods := restHandler.GetMethods()
	for _, method := range methods {
		this.PutMethod(method, restHandler.GetSigninHandler())
	}
}
