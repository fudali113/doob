package router

import (
	"log"
	"testing"
)

func Test_createNodeValue(t *testing.T) {
	testStr := "{ooo}-{hahaha:[0-9]{3}}"
	v := createNodeValue(testStr)

	log.Println("Test_createNodeValue : ", v)
}
