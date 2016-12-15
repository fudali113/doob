package return_deal

import (
	"encoding/json"
	"net/http"
)

//	处理返回 type 为 json 或者返回单个对象的 func
//	返回单个对象且 type 不为 string 的 handle func 默认返回处理 type 为 json
type ReturnJsonDealer struct {
	DealerName
}

//	实现 Deal 方法
func (*ReturnJsonDealer) Deal(returnType *ReturnType, w http.ResponseWriter) {
	var data interface{}
	if returnType.Data == nil {
		data = map[string]string{}
	} else {
		data = returnType.Data
	}
	json, err := json.Marshal(data)
	if err != nil {
		// TODO log
		return
	}
	w.Write(json)
	w.Header().Add("context", "application/json")
}

func (*ReturnJsonDealer) MacthType(str string) bool {
	return matchPrefix(str, "json")
}

func init() {
	AddReturnDealer(&ReturnJsonDealer{DealerName: DealerName{name: DEFAULT_JSON_DEALER_NAME}})
}
