package router

import (
	"log"
	"testing"
)

var (
	testData = []string{
		"/admin/login",
		"/admin/system/log/return_value/{on_off}",
		"/api/user/barcodes",
		"/api/user/barcodes/{barcode}",
		"/api/user/barcodes/{barcode}/auth/{toUserId}",
		"/api/user/barcodes/{barcode}/share",
		"/api/user/barcodes/{barcode}/bind_share",
		"/api/user/info",
		"/api/user/password",
		"/api/user/coupon",
		"/api/user/coupon/{couponId}",
		"/api/user/feedback",
		"/api/login",
		"/api/sign_in",
		"api/report/cancer/{id}",
		"api/report/cancer/static/{id}",
		"api/report/cancer/gene/{unionId}",
		"api/report/cancer/gene/static/{unionId}",
		"api/report/disease/",
		"api/report/disease/static/{id}",
		"api/report/disease/{id}",
		"/api/report/drug/",
		"/api/report/drug/class",
		"/api/report/drug/class/{name}",
		"/api/report/drug/{id}",
		"api/report/inherit/",
		"api/report/inherit/static/{id}",
		"api/report/inherit/{id}",
		"api/report/nutrition/",
		"api/report/nutrition/{id}",
		"api/report/nutrition/static/{id}",
		"/api/report/",
		"/api/report/index",
		"/api/report/all_item",
		"/api/report/more_tags",
		"api/report/sport/",
		"api/report/sport/{id}",
		"api/report/sport/tatic/{id}",
		"api/report/sport/rare_gene/{id}",
		"/api/report/trait",
		"/api/report/trait/class",
		"/api/report/trait/class/{id}",
		"/api/report/trait/{id}",
		"/api/report/trait/static/{id}",
		"/api/report/trait/rare_gene/{id}",
	}
)

type testType struct {
	name string
	num  int
}

func (t testType) String() string {
	return t.name + "----"
}

func Benchmark_test(b *testing.B) {
	node := &Node{
		Class:    normal,
		Value:    nil,
		handler:  nil,
		Children: new(childrens),
	}
	testVar := &RegisterHandler{
		Handler: &testType{
			name: "ooo",
			num:  1,
		},
	}
	for _, url := range testData {
		node.InsertChild(url, GetSimpleRestHandler("get", testVar))
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res := must(node.GetRT("/api/report/index", nil)).GetHandler("get").GetHandler()
			if res == nil {
				b.Error("path variable method have bug")
			}

			res1 := must(node.GetRT("/api/user/barcodes/111-1121-8406/bind_share", nil)).GetHandler("get").GetHandler()
			if res1 == nil {
				b.Error("path variable method have bug")
			}
		}
	})
}

func must(i RestHandler, e error) RestHandler {
	if e != nil {
		log.Panic(e)
	}
	return i
}
