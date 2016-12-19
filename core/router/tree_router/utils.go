package tree_router

import "strings"

func getUrlNodeValue(url string) (string ,  string) {
	url = strings.TrimPrefix(url,"/")
	prefixAndSuffix := strings.SplitN(url,"/",1)
	return prefixAndSuffix[0],prefixAndSuffix[1]
}
