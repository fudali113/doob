package return_deal

import (
	"html/template"
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

// 实现 Dealer 接口
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
		getTemplateBytes(path, data, w)
	}
}

// 获取 typeStr 中的路径
func getPath(typeStr string) string {
	typeAndPath := utils.Split(typeStr, ":")
	if len(typeAndPath) == 2 {
		return typeAndPath[1]
	}
	return ""
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
	AddReturnDeal(&staticHtmlReturnDeal{})
}
