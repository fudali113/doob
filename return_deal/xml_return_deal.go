package return_deal

import (
	"encoding/xml"
	"log"
	"net/http"

	. "github.com/fudali113/doob/http_const"
)

//	处理返回 type 为 xml 或者返回单个对象的 func
type ReturnXmlDealer struct {
	DealerName
}

func (*ReturnXmlDealer) MacthType(str string) bool {
	return matchPrefix(str, "xml")
}

//	实现 Deal 方法
func (*ReturnXmlDealer) Deal(returnType *ReturnType, w http.ResponseWriter) {
	var data interface{}
	if returnType.Data == nil {
		data = map[string]string{}
	} else {
		data = returnType.Data
	}
	xml, err := xml.Marshal(data)
	if err != nil {
		log.Print("xml dealer is error , error is ", err)
		return
	}
	w.Header().Add(CONTENT_TYPE, APP_XML)
	w.Write(xml)
}

func init() {
	AddReturnDealer(&ReturnXmlDealer{DealerName: DealerName{name: DEFAULT_XML_DEALER_NAME}})
}
