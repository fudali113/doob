package register

type RegisterHandlerType interface {
	GetHandler() interface{}
	GetRegisterType() *RegisterType
}
