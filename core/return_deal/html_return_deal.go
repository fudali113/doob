package return_deal

import (
	"io/ioutil"
	"log"
	"net/http"
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
		w.WriteHeader(500)
		return
	}
	w.Write(bytes)
}

func init() {
	AddReturnDealer(&staticFileReturnDealer{DealerName: DealerName{name: DEFAULT_HTML_DEALER_NAME}})
}
