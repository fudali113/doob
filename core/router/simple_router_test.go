package router

import (
	"log"
	"testing"
)

type testType struct {
	name string
	num  int
}

func Test_getAndAdd(t *testing.T) {
	simpleRouter := &SimpleRouter{}
	testVar := &testType{
		name: "ooo",
		num:  1,
	}
	testVar1 := &testType{
		name: "ooo222",
		num:  2,
	}
	/**
	 * normal
	 */
	simpleRouter.Add("/dddd/dssds/dfdggf", &SimpleRestHandler{
		signinHandler: testVar,
		methodHandlerMap: map[string]interface{}{
			"get": testVar,
		},
	})
	simpleRouter.Add("/dddd/dssds/dfdggf", &SimpleRestHandler{
		signinHandler: testVar1,
		methodHandlerMap: map[string]interface{}{
			"post": testVar1,
		},
	})
	testGetVar, _ := simpleRouter.Get("/dddd/dssds/dfdggf").Rest.GetHandler("get").(*testType)
	_testGetVar, _ := simpleRouter.Get("/dddd/dssds/dfdggf").Rest.GetHandler("post").(*testType)
	if testGetVar != testVar {
		t.Error("normal router is error")
	}
	if testGetVar.num != 1 {
		t.Error("normal router is error")
	}
	if _testGetVar != testVar1 {
		t.Error("normal router is error")
	}
	if _testGetVar.num != 2 {
		t.Error("normal router is error")
	}

	/**
	 * pathVariable
	 */
	simpleRouter.Add("/dddd/{fffff}/dfdggf", &SimpleRestHandler{
		signinHandler: testVar,
		methodHandlerMap: map[string]interface{}{
			"get": testVar,
		},
	})
	simpleRouter.Add("/dddd/{fffff}/dfdggf", &SimpleRestHandler{
		signinHandler: testVar1,
		methodHandlerMap: map[string]interface{}{
			"post": testVar1,
		},
	})
	for i := 0; i < 10; i++ {
		testGetVar1, _ := simpleRouter.Get("/dddd/dssssssds/dfdggf").Rest.GetHandler("get").(*testType)
		_testGetVar1, _ := simpleRouter.Get("/dddd/dssssssds/dfdggf").Rest.GetHandler("post").(*testType)
		if testGetVar1 != testVar {
			t.Error("normal router is error")
		}
		if testGetVar1.num != 1 {
			t.Error("normal router is error")
		}
		if _testGetVar1 != testVar1 {
			t.Error("normal router is error")
		}
		if _testGetVar1.num != 2 {
			t.Error("normal router is error")
		}
	}

	/**
	 * suffix
	 */
	simpleRouter.Add("/ddf/**", &SimpleRestHandler{
		signinHandler: testVar,
		methodHandlerMap: map[string]interface{}{
			"get": testVar,
		},
	})
	testGetVar2, _ := simpleRouter.Get("/ddf/dssds/dfdggf,dsds-+!@#$%^&*").Rest.GetHandler("get").(*testType)
	if testGetVar2 != testVar {
		t.Error("normal router is error")
	}
	if testGetVar2.num != 1 {
		t.Error("normal router is error")
	}

}

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

/**
 * 测试获取是否pathVariable匹配url的正则
 */
func Test_getPathVariableReg(t *testing.T) {
	testUrlStr := "/api/dd/{ddddd}/{ffffffff}/fdldfldf"
	shouldRegStr := "/api/dd/\\S+/\\S+/fdldfldf"
	oo := getPathVariableReg(testUrlStr)
	if shouldRegStr != oo.String() {
		t.Error("getPathVariableReg func have a bug")
	}
}

func Test_getPathVariableParamMap(t *testing.T) {
	pathVariableHandler := getPathVariableHandler("/{name1}/dfdfdf_{name2}/dfdf_{name3}_dsffdffd/{name4}", nil)
	res := pathVariableHandler.getPathVariableParamMap("/name1/dfdfdf_name2/dfdf_name3_dsffdffd/name4")
	if res["name1"] != "name1" {
		t.Error("getPathVariableParamMap func have bug")
	}
	if res["name2"] != "name2" {
		t.Error("getPathVariableParamMap func have bug")
	}
	if res["name3"] != "name3" {
		t.Error("getPathVariableParamMap func have bug")
	}
	if res["name4"] != "name4" {
		t.Error("getPathVariableParamMap func have bug")
	}
}

func Test_getRegexp(t *testing.T) {
	testStr := "{sddssddssd}dfdfdf{ooooooo}dfdf{ooooooooo}dsffdffd"
	r := getRegexp("{\\w+}")
	// 	s := r.Split(testStr, -1)
	// 	log.Printf("FindString res:%s", s)
	//
	s1 := r.FindAllStringSubmatch(testStr, -1)
	log.Printf("FindAllStringSubmatch res:%s", s1)
	// 	s2 := r.FindStringIndex(testStr)
	// 	log.Printf("FindAllStringSubmatch res:%d", s2)
	// 	log.Print("FindStringSubmatch : ", r.FindStringSubmatch(testStr))
}
