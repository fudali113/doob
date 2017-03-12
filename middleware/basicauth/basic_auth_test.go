package basicauth

import "testing"

func Test_getUsernameAndPasswd(t *testing.T) {
	testStr := "YWRtaW46MTIzNDU2"
	username, passwd, err := getUsernameAndPasswd(testStr)
	if err != nil || username != "admin" || passwd != "123456" {
		t.Error(" func  getUsernameAndPasswd has bug")
	}
}
