package doob

var (
	urlSplitSymbol         = "&&"
	returnDealDefaultType  = "auto"
	redirectDefaulBodytLen = 1024
)

func SetReturnDealDefaultType(t string) {
	returnDealDefaultType = t
}

func SetUrlSplitSymbol(symbol string) {
	urlSplitSymbol = symbol
}

func SetRedirectDefaulBodytLen(len int) {
	redirectDefaulBodytLen = len
}
