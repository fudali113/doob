package utils

import "testing"

func Test_Split(t *testing.T) {

	a := "  fgfdgf     ffhbgf    fdfbfg   "
	b := Split(a, " ")
	//AssertLT(t,3,2)
	AssertEQ(t, len(b), 3)

}

func Test_GetGroupMatch(t *testing.T) {
	test := "aaaa{aaa}-{bbbb:[abc]{3}}dddddd"
	m, nm := GetGroupMatch(test, '{', '}')
	if m[0] != "{aaa}" || m[1] != "{bbbb:[abc]{3}}" || nm[0] != "aaaa" {
		t.Error("GetGroupMatch has bug")
	}
}
