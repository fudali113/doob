package golib

import "testing"


func AssertEQ(t *testing.T , fact ,hope interface{})  {
	if hope != fact{
		t.Errorf(" false ->  %v == %v ",hope,fact)
	}
}

func AssertGT(t *testing.T , fact ,hope int){
	if hope >= fact{
		t.Errorf(" false ->  %v > %v " ,fact ,hope)
	}
}

func AssertGE(t *testing.T , fact ,hope int){
	if hope > fact{
		t.Errorf(" false ->  %v >= %v " ,fact ,hope)
	}
}

func AssertLT(t *testing.T , fact ,hope int){
	if hope <= fact{
		t.Errorf(" false ->  %v < %v " ,fact ,hope)
	}
}

func AssertLE(t *testing.T , fact ,hope int){
	if hope < fact{
		t.Errorf(" false ->  %v <= %v " ,fact ,hope)
	}
}

