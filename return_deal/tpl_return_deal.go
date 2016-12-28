package return_deal

import (
	"html/template"
	"log"
	"net/http"

	. "github.com/fudali113/doob/http_const"
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
	data := returnType.Data
	if data != nil {
		getTemplateBytes(path, data, w)
	}
}

// 标准 template 实现
func getTemplateBytes(path string, data interface{}, w http.ResponseWriter) {
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Print("tpl dealer is error , error is ", err)
		w.WriteHeader(INTERNAL_SERVER_ERROR)
		return
	}
	w.Header().Add(CONTENT_TYPE, APP_HTML)
	err = t.Execute(w, data)
	if err != nil {
		log.Print("tpl dealer is error , error is ", err)
		w.WriteHeader(INTERNAL_SERVER_ERROR)
	}
}

func init() {
	AddReturnDealer(&tplReturnDeal{DealerName: DealerName{name: DEFAULT_TPL_DEALER_NAME}})
}
