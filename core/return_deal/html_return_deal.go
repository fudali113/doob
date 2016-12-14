package return_deal

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fudali113/doob/utils"
)

type staticFileReturnDealer struct {
}

func (*staticFileReturnDealer) MacthType(str string) bool {
	return strings.HasPrefix(str, "html") || strings.HasPrefix(str, "file")
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

// 获取 typeStr 中的路径
func getPath(typeStr string) string {
	typeAndPath := utils.Split(typeStr, ":")
	if len(typeAndPath) == 2 {
		path := typeAndPath[1]
		suffix := "." + typeAndPath[0]
		// 如果文件没有跟上后缀名
		// 加上后缀名
		if !strings.HasSuffix(path, suffix) {
			path = path + suffix
		}
		return path
	}
	return ""
}

func init() {
	AddReturnDealer(&staticFileReturnDealer{})
}
