package return_deal

import (
	"encoding/json"
	"log"
	"net/http"

	. "github.com/fudali113/doob/core/http_const"
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
		log.Print("json dealer is error , error is ", err)
		return
	}
	w.Header().Add(CONTENT_TYPE, APP_JSON)
	w.Write(json)
}

func (*ReturnJsonDealer) MacthType(str string) bool {
	return matchPrefix(str, "json")
}

func init() {
	AddReturnDealer(&ReturnJsonDealer{DealerName: DealerName{name: DEFAULT_JSON_DEALER_NAME}})
}
