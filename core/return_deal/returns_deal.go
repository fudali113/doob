package return_deal

import "net/http"

var (
	deals = make([]ReturnMatchType, 0)
)

/**
 * 处理
 */
func DealReturn(returnType ReturnType, req http.ResponseWriter) {
	for _, returnDeal := range deals {
		if returnDeal.MacthType(returnType.typeStr) {
			serialize, ok := returnDeal.(Serialize)
			if ok {
				bytes, headers := serialize.Serialize(returnType)
				req.Write(bytes)
				for name, value := range map[string][]string(headers) {
					req.Header().Add(name, value[0])
				}
				return
			}
			deal, ok := returnDeal.(Deal)
			if ok {
				deal.Deal(returnType, req)
				return
			}

		}
	}
}

/**
 * 添加一个处理实例
 */
func AddReturnDeal(returnDeals ...ReturnMatchType) {
	deals = append(deals, returnDeals...)
}
