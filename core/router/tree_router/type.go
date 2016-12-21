package tree_router

import (
	"log"
	"regexp"
	"strings"
)

// 各类型储存接口
type nodeV interface {
	isMatch(urlPart string) (bool, bool)
	// if need pathvar
	// return in this method
	paramValue(urlPart string, url string) map[string]string
	getOrigin() string
}

type nodeVNormal struct {
	origin string
}

func (this nodeVNormal) isMatch(urlPart string) (bool, bool) {
	return this.origin == urlPart, false
}
func (this nodeVNormal) paramValue(urlPart string, url string) map[string]string {
	return map[string]string{}
}
func (this nodeVNormal) getOrigin() string {
	return this.origin
}

type nodeVPathReg struct {
	origin    string
	paramName string
	paramReg  *regexp.Regexp
}

// check url part is match this node value
func (this nodeVPathReg) isMatch(urlPart string) (bool, bool) {
	findStr := this.paramReg.FindString(urlPart)
	log.Print(findStr, "====", urlPart)
	return findStr == urlPart, false
}
func (this nodeVPathReg) paramValue(urlPart string, url string) map[string]string {
	return map[string]string{this.paramName: urlPart}
}
func (this nodeVPathReg) getOrigin() string {
	return this.origin
}

type nodeVPathVar struct {
	origin    string
	paramName string
}

func (this nodeVPathVar) isMatch(urlPart string) (bool, bool) {
	return true, false
}
func (this nodeVPathVar) paramValue(urlPart string, url string) map[string]string {
	return map[string]string{this.paramName: urlPart}
}
func (this nodeVPathVar) getOrigin() string {
	return this.origin
}

type nodeVMatchAll struct {
	origin string
	prefix string
}

func (this nodeVMatchAll) isMatch(urlPart string) (bool, bool) {
	return strings.HasPrefix(urlPart, this.prefix), true
}
func (this nodeVMatchAll) paramValue(urlPart string, url string) map[string]string {
	return map[string]string{"**": strings.TrimPrefix(url, this.prefix)}
}
func (this nodeVMatchAll) getOrigin() string {
	return this.origin
}
