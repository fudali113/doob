package doob

import (
	"testing"
	"net/url"
)

func Test_urlParse(t *testing.T)  {
	testUrl := "/ooooo"
	url , _ := url.Parse(testUrl)
	t.Log(url.Host, url.Path)
	testUrl1 := "https://github.com/astaxie/beego/blob/master/context/context.go"
	url1 , _ := url.Parse(testUrl1)
	t.Log(url1.Host, url1.Path, url1.Scheme)
}
