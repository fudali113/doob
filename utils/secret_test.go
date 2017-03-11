package utils

import (
	"log"
	"testing"

	"github.com/fudali113/doob/config"
)

func Test_GetRandomStr(t *testing.T) {
	randStr := GetRandomStr()
	if randStr == "" {
		t.Error("GetRandomStr func has bug")
	}
}

func Test_GetMd5String(t *testing.T) {
	randStr := GetRandomStr()
	md5Str := GetMd5String(randStr, config.SessionCreateSecretKey)
	log.Print(md5Str, "---", len(md5Str))
	if len(md5Str) != 32 {
		t.Error("GetMd5String func has bug")
	}
}
