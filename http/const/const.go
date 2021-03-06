//
// http related const store
//
package http

const (
	GET     = "get"
	POST    = "post"
	PUT     = "put"
	DELETE  = "delete"
	HEAD    = "head"
	OPTIONS = "options"

	CACHE_CONTROL = "Cache-Control"
	NO_CACHE      = "no-cache"
	ACCEPT        = "Accept"
	LOCATION      = "Location"
	SET_COOKIE    = "Set-Cookie"
	COOKIE        = "Cookie"
	BASIC_AUTH    = "Authorization"
	WWW_AUTH      = "WWW-Authenticate"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Allow
	ALLOW_METHODS = "Allow"
	SERVER        = "Server"

	CONTENT_TYPE = "Content-Type"
	APP_JSON     = "application/json"
	APP_XML      = "application/xml"
	APP_HTML     = "text/html"
	MULTIPART    = "multipart/form-data"
)
