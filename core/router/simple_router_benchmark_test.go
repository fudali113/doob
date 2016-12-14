package router

import (
	"testing"

	"github.com/fudali113/doob/core/register"
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

func Benchmark_test(b *testing.B) {
	simpleRouter := &SimpleRouter{}
	testVar := &RegisterHandler{
		Handler: &testType{
			name: "ooo",
			num:  1,
		},
	}
	for _, url := range testData {
		simpleRouter.Add(url, &SimpleRestHandler{
			methodHandlerMap: map[string]register.RegisterHandlerType{
				"get": testVar,
			},
		})
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			oo, _ := simpleRouter.Get("/api/report/index" /**fmt.Sprintf("/aaa/%s/oo", "dd")*/).Rest.GetHandler("get").GetHandler().(*testType)
			if oo == nil {
				b.Error("path variable method have bug")
			}

			oo1, _ := simpleRouter.Get("/api/user/barcodes/111-1121-8406/bind_share").Rest.GetHandler("get").GetHandler().(*testType)
			if oo1 == nil {
				b.Error("path variable method have bug")
			}
		}
	})
}
