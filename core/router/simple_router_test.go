package router

import (
	"log"
	"testing"
)

/**
 * 测试url类型，分类处理
 */
func Test_getUrlClassify(t *testing.T) {
	res1 := getUrlClassify("/oooo/dddd/ddddd/**")
	if res1 != LAST_ALL_MATCH {
		t.Errorf("should 2 but %d", res1)
	}
	res2 := getUrlClassify("/dddd/eeee{ssdsdsd}/**")
	if res2 != PV_AND_LAM {
		t.Errorf("should 3 but %d", res2)
	}
	res3 := getUrlClassify("/dddd/eeee/dddd")
	if res3 != NORMAL {
		t.Errorf("should 0 but %d", res3)
	}
	res4 := getUrlClassify("/fsdfdsf/sfdsf{dsdssd}/dsfsdsf")
	if res4 != PATH_VARIABLE {
		t.Errorf("should 1 but %d", res3)
	}
}

func Test_getPathVariableReg(t *testing.T) {
	log.Print(getPathVariableReg("/api/dd/{ddddd}/{ffffffff}/fdldfldf"))
}

// func Test_getRegexp(t *testing.T) {
// 	testStr := "{sddssddssd}dfdfdfdfdf"
// 	r := getRegexp("{\\w+}")
// 	log.Print(r)
// 	s := r.FindString(testStr)
// 	log.Print(s == "")
// 	log.Printf("FindString res:%s", s)
// 	s1 := r.FindAllStringSubmatch(testStr, 1)
// 	log.Printf("FindAllStringSubmatch res:%s", s1)
// 	s2 := r.FindStringIndex(testStr)
// 	log.Printf("FindAllStringSubmatch res:%d", s2)
// 	log.Print("FindStringSubmatch : ", r.FindStringSubmatch(testStr))
// 	t.Log(s)
// }
