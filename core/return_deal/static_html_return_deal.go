package return_deal

import (
	"log"
	"net/http"
	"strings"
)

type staticHtmlReturnDeal struct {
}

func (*staticHtmlReturnDeal) MacthType(str string) bool {
	return strings.HasPrefix(str, "html")
}

func (*staticHtmlReturnDeal) Deal(returnType ReturnType, res http.ResponseWriter) {
	log.Print(returnType.TypeStr)
}

func init() {
	AddReturnDeal(&staticHtmlReturnDeal{})
}
