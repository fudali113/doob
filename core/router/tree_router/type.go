package tree_router

import "regexp"

type nodeV interface {
	isMatch(urlPart string) bool
	paramValue(urlPart string) map[string]string
}

type nodeValue struct {
	class int
	origin string
	paramName string
	paramReg *regexp.Regexp
}

// check url part is match this node value
func (this *nodeValue) isMatch(urlPart string) bool {
	switch this.class {
	case normal:
		return this.origin == urlPart
	case pathReg:
		return this.paramReg.MatchString(urlPart)
	case pathVar,matchAll:
		return true
	}
	return false
}
