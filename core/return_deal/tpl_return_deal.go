package return_deal

import (
	"html/template"
	"log"
	"net/http"
)

type tplReturnDeal struct {
	DealerName
}

func (*tplReturnDeal) MacthType(str string) bool {
	return matchPrefix(str, "tpl")
}

// 实现 Dealer 接口
func (*tplReturnDeal) Deal(returnType *ReturnType, w http.ResponseWriter) {
	path := getPath(returnType.TypeStr)
	log.Print(path)
	data := returnType.Data
	if data != nil {
		getTemplateBytes(path, data, w)
	}
}

// 标准 template 实现
func getTemplateBytes(path string, data interface{}, w http.ResponseWriter) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(500)
	}
}

func init() {
	AddReturnDealer(&tplReturnDeal{DealerName: DealerName{name: DEFAULT_HTML_DEALER_NAME}})
}
