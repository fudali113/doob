package return_deal

import "net/http"

type ReturnType struct {
	TypeStr string
	Data    interface{}
}

type ReturnMatchType interface {
	/**
	 * 是否匹配typeStr
	 */
	MacthType(typeStr string) bool
}

type Serialize interface {
	/**
	 * 序列化方式，
	 */
	Serialize(returnType *ReturnType) ([]byte, http.Header)
}

type Deal interface {
	/**
	 * 自己处理相关数据到respons
	 */
	Deal(returnType *ReturnType, res http.ResponseWriter)
}
