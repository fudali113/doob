package utils

import (
	"log"
	"regexp"
	"strings"
)

// Split golang1.6 Strings.split()有问题
// 如果以``" "`分割将出现很多`""`字符串
func Split(s, sep string) []string {
	sepSave := 0
	n := strings.Count(s, sep) + 1
	c := sep[0]
	start := 0
	a := make([]string, n)
	na := 0
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		if s[i] == c && (len(sep) == 1 || s[i:i+len(sep)] == sep) {
			splitStr := s[start : i+sepSave]
			if !(splitStr == sep || start == i+sepSave) {
				a[na] = splitStr
				na++
			}
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	if last := s[start:]; last != "" {
		a[na] = last
		na++
	}
	return a[:na]
}

// GetRegexp 获取正则表达式
func GetRegexp(reg string) *regexp.Regexp {
	r, err := regexp.Compile(reg)
	if err != nil {
		log.Print(err)
		return nil
	}
	return r
}

// GetGroupMatch 成对的获取字符串
// md ， 正则平衡组太难理解
func GetGroupMatch(origin string, prefix, suffix rune) (match, notMatch []string) {
	originRunes := []rune(origin)
	getGroupMatch(originRunes, prefix, suffix, &match, &notMatch)
	return
}

func getGroupMatch(originRunes []rune, prefix, suffix rune, match, notMatch *[]string) {
	prefixMatchIndex := 0
	suffixMatchIndex := 0
	nowPrefixMatchNum := 0
	nowSuffixMatchNum := 0
	for i := 0; i < len(originRunes); i++ {
		nowRune := originRunes[i]
		if nowRune == prefix {
			if nowPrefixMatchNum == 0 {
				prefixMatchIndex = i
			}
			nowPrefixMatchNum++
		}
		if nowRune == suffix && nowPrefixMatchNum > nowSuffixMatchNum {
			suffixMatchIndex = i
			nowSuffixMatchNum++
		}
		if nowPrefixMatchNum > 0 && nowPrefixMatchNum == nowSuffixMatchNum {
			*notMatch = append(*notMatch, string(originRunes[0:prefixMatchIndex]))
			*match = append(*match, string(originRunes[prefixMatchIndex:suffixMatchIndex+1]))
			nextOrigin := originRunes[suffixMatchIndex+1:]
			if len(nextOrigin) > 0 {
				getGroupMatch(nextOrigin, prefix, suffix, match, notMatch)
			}
			break
		}
		if i == len(originRunes)-1 {
			*notMatch = append(*notMatch, string(originRunes))
		}
	}
}
