package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fudali113/doob/utils"

	. "github.com/fudali113/doob/http/const"
)

const (
	debugHTMLTplStr = `
		<html>
			<body>
				<font size="8" color="red">error info : </font>
				<div>%v</div>
				<font size="8" color="red">stack info : </font>
				<div>%s</div>
			<body>
		</html>
	`
)

// WriteErrInfo 返回错误信息
func WriteErrInfo(err interface{}, stack []byte, w http.ResponseWriter) {
	stackStr := string(stack)
	stackSlice := utils.Split(stackStr, "\n")
	stackStr = ""
	panic := false
	for _, str := range stackSlice {
		html := ""
		isInfo := strings.HasPrefix(str, "\t")
		if strings.HasPrefix(str, "panic") {
			panic = true
		} else if !isInfo {
			panic = false
		}

		background := ""
		if panic {
			background = "background:red"
		}

		if isInfo {
			html = `<div style="padding-left:31;` + background + `">` + str + "</div>"
		} else {
			html = `<div style="` + background + `">` + str + "</div>"
		}
		stackStr += html
	}
	html := fmt.Sprintf(debugHTMLTplStr, err, stackStr)
	w.Header().Add(CONTENT_TYPE, APP_HTML)
	w.Write([]byte(html))
}
