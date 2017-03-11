package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// 以随机数为key生成32为随机数
func GetRandomStr() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return string(b)
}

//生成32位md5字串
func GetMd5String(s, scretKey string) string {
	h := md5.New()
	h.Write([]byte(s + scretKey))
	return hex.EncodeToString(h.Sum(nil))
}
