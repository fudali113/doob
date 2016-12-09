package return_deal

import "net/http"

var (
	deals = make([]ReturnMatchType, 0)
)

/**
 * 处理
 */
func DealReturn(returnType *ReturnType, w http.ResponseWriter) {
	for _, returnDeal := range deals {
		if returnDeal.MacthType(returnType.TypeStr) {
			serialize, ok := returnDeal.(Serialize)
			if ok {
				bytes, headers := serialize.Serialize(returnType)
				w.Write(bytes)
				for name, value := range map[string][]string(headers) {
					w.Header().Add(name, value[0])
				}
				return
			}
			deal, ok := returnDeal.(Deal)
			if ok {
				deal.Deal(returnType, w)
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
