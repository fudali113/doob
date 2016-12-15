package return_deal

import "net/http"

type ReturnType struct {
	TypeStr string
	Data    interface{}
}

type ReturnTypeDealer interface {
	// 此 Dealer 的名字
	Name() string
	//	是否匹配typeStr
	MacthType(typeStr string) bool
	//	自己处理相关数据到respons
	Deal(returnType *ReturnType, w http.ResponseWriter)
}

type DealerName struct {
	name string
}

func (this DealerName) Name() string {
	return this.name
}
