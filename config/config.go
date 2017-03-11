package config

var (
	IsDev                  = true
	AutoAddHead            = true
	AutoAddOptions         = true
	UrlSplitSymbol         = "&&"
	ReturnDealDefaultType  = "auto"
	RedirectDefaultBodyLen = 1024

	BasicAuthReminder = "doob basic auth"

	SessionCreateSecretKey = "doobssssssss"
	OpenBasicAuth          = false
	BasicAuthUserInConfig  = map[string]string{"admin": "123456"}
)
