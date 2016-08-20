package golib

import (
	"testing"
)

func Test_Split(t *testing.T){

	a:="  fgfdgf     ffhbgf    fdfbfg   "
	b:=Split(a," ")
	AssertLT(t,3,2)
	AssertEQ(t,len(b),3)

}
