package doob

import (
	"net/url"
	"testing"
)

func Test_urlParse(t *testing.T) {
	testURL := "/ooooo"
	url, _ := url.Parse(testURL)
	t.Log(url.Host, url.Path)
	testURL1 := "https://github.com/astaxie/beego/blob/master/context/context.go"
	url1, _ := url.Parse(testURL1)
	t.Log(url1.Host, url1.Path, url1.Scheme)
}
