package return_deal

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ReturnJsonSerialize struct {
}

func (*ReturnJsonSerialize) Serialize(returnType ReturnType) ([]byte, http.Header) {
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
