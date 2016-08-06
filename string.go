package golib

import "strings"
func Split(s, sep string) []string {
	sepSave := 0
	n := strings.Count(s, sep) + 1
	c := sep[0]
	start := 0
	a := make([]string, n)
	na := 0
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		if s[i] == c && (len(sep) == 1 || s[i:i+len(sep)] == sep) {
			splitStr:=s[start : i+sepSave]
			if !(splitStr == sep || start==i+sepSave) {
				a[na] = splitStr
				na++
			}
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	if last:=s[start:];last!=""{a[na] = last}
	return a[: na]
}