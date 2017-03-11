package config

var (
	IsDev                  = true
	AutoAddHead            = true
	AutoAddOptions         = true
	UrlSplitSymbol         = "&&"
	ReturnDealDefaultType  = "auto"
	RedirectDefaultBodyLen = 1024

	SessionCreateSecretKey = "doobssssssss"
	OpenBasicAuth          = true
	BasicAuthUserInConfig  = map[string]string{"admin": "123456"}
)
