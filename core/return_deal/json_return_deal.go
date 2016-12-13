package return_deal

import (
	"encoding/json"
	"net/http"
	"strings"
)

//	处理返回 type 为 json 或者返回单个对象的 func
//	返回单个对象且 type 不为 string 的 handle func 默认返回处理 type 为 json
type ReturnJsonSerialize struct {
}

//	实现 Serializer 接口
func (*ReturnJsonSerialize) Serialize(returnType *ReturnType) ([]byte, http.Header) {
	header := http.Header{}
	var data interface{}
	if returnType.Data == nil {
		data = map[string]string{}
	} else {
		data = returnType.Data
	}
	json, err := json.Marshal(data)
	if err != nil {

	}
	header.Add("json", "application/json")
	return json, header
}

func (*ReturnJsonSerialize) MacthType(str string) bool {
	return strings.ToLower(str) == "json"
}

func init() {
	AddReturnDeal(&ReturnJsonSerialize{})
}
