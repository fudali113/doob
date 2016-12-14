package return_deal

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fudali113/doob/utils"
)

type staticHtmlReturnDeal struct {
}

func (*staticHtmlReturnDeal) MacthType(str string) bool {
	return strings.HasPrefix(str, "html")
}

//	实现 Dealer 接口
func (*staticHtmlReturnDeal) Deal(returnType *ReturnType, w http.ResponseWriter) {
	log.Print(returnType.TypeStr)
	path := getPath(returnType.TypeStr)
	data := returnType.Data
	if data == nil {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(bytes)
	} else {
		bytes, err := getTempleteBytes(path, data)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(bytes)
	}
}

func getPath(typeStr string) string {
	typeAndPath := utils.Split(typeStr, ":")
	if len(typeAndPath) == 2 {
		return typeAndPath[1]
	}
	return ""
}

func getTempleteBytes(path string, data interface{}) ([]byte, error) {
	return []byte{}, nil
}

func init() {
	AddReturnDeal(&staticHtmlReturnDeal{})
}
