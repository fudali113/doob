package golib

import (
	"testing"
	"fmt"
)

func Test_Split(t *testing.T){

	a:="  fgfdgf     ffhbgf    fdfbfg   "
	b:=Split(a," ")

	if len(b)!=3{
		t.Error(fmt.Sprintf("Split() is error! %d",len(b)))
	}

}
