package return_deal

import (
	"io/ioutil"
	"log"
	"net/http"

	. "github.com/fudali113/doob/http/const"
)

type staticFileReturnDealer struct {
	DealerName
}

func (*staticFileReturnDealer) MacthType(str string) bool {
	return matchPrefix(str, "html", "file")
}

// 实现 Dealer 接口
func (*staticFileReturnDealer) Deal(returnType *ReturnType, w http.ResponseWriter) {
	log.Print(returnType.TypeStr)
	path := getPath(returnType.TypeStr)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print("html/file dealer is error , error is ", err)
		w.WriteHeader(NOT_FOUND)
		return
	}
	w.Header().Add(CONTENT_TYPE, APP_HTML)
	w.Write(bytes)
}

func init() {
	AddReturnDealer(&staticFileReturnDealer{DealerName: DealerName{name: DEFAULT_HTML_DEALER_NAME}})
}
