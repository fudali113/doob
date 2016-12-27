//
//	处理用户注册的 handle func 的返回值
//	列如
//
//		func (c *core.Context) (string , inteface{}){}
//
//	它的 type 是返回的string类型的值
//	处理的 data 是inteface{} 类型的返回值
//
//	使用者在应用初始化时添加相关的处理函数
//	如有 handler func 的返回值匹配用户实现 ReturnMatchType 接口的 MacthType 方法时
//	将进入到用户实现的 Serializer 或者 Dealer 接口的方法中并返回
//
package return_deal

import (
	"log"
	"net/http"
)

// default dealer names
// user can rewriter this name dealer
// cover default dealer
const (
	DEFAULT_JSON_DEALER_NAME = "json_dealer"
	DEFAULT_XML_DEALER_NAME  = "xml_dealer"
	DEFAULT_HTML_DEALER_NAME = "html_dealer"
	DEFAULT_TPL_DEALER_NAME  = "tpl_dealer"
)

// dealer store
var (
	dealerMap = make(map[string]ReturnTypeDealer, 0)
)

// 根据初始化时加入的元素进行遍历处理
// 找到第一个匹配的type是进行处理
// 之后将结束遍历并返回
func DealReturn(returnType *ReturnType, w http.ResponseWriter) {
	for _, dealer := range dealerMap {
		if dealer.MacthType(returnType.TypeStr) {
			dealer.Deal(returnType, w)
			return
		}
	}
	log.Print("don`t have deal handler match this type : ", returnType.TypeStr)
}

//	添加一个处理实例
func AddReturnDealer(returnDeals ...ReturnTypeDealer) {
	for _, returnDeal := range returnDeals {
		if _, ok := dealerMap[returnDeal.Name()]; ok {
			log.Print("has repetition dealer , name is ", returnDeal.Name())
		}
		dealerMap[returnDeal.Name()] = returnDeal
	}
}
