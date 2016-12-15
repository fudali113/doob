package return_deal

import (
	"encoding/xml"
	"net/http"
)

//	处理返回 type 为 xml 或者返回单个对象的 func
type ReturnXmlDealer struct {
	DealerName
}

func (*ReturnXmlDealer) MacthType(str string) bool {
	return matchPrefix(str, "xml")
}

//	实现 Serializer 接口
func (*ReturnXmlDealer) Deal(returnType *ReturnType, w http.ResponseWriter) {
	var data interface{}
	if returnType.Data == nil {
		data = map[string]string{}
	} else {
		data = returnType.Data
	}
	xml, err := xml.Marshal(data)
	if err != nil {
		// TODO log
		return
	}
	w.Write(xml)
	w.Header().Add("context", "application/xml")
}

func init() {
	AddReturnDealer(&ReturnXmlDealer{DealerName: DealerName{name: DEFAULT_XML_DEALER_NAME}})
}
