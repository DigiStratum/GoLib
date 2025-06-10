package http

import "strings"

type HttpRequestMethod int

const (
	METHOD_UNKNOWN HttpRequestMethod = iota
	METHOD_GET
	METHOD_POST
	METHOD_DELETE
	METHOD_PATCH
	METHOD_PUT
	METHOD_HEAD
	METHOD_OPTIONS
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func HttpRequestMethodFromString(httpRequestMethod string) HttpRequestMethod {
	switch strings.ToUpper(httpRequestMethod) {
	case "GET":
		return METHOD_GET
	case "POST":
		return METHOD_POST
	case "DELETE":
		return METHOD_DELETE
	case "PATCH":
		return METHOD_PATCH
	case "PUT":
		return METHOD_PUT
	case "HEAD":
		return METHOD_HEAD
	case "OPTIONS":
		return METHOD_OPTIONS
	}
	return METHOD_UNKNOWN
}

// -------------------------------------------------------------------------------------------------
// HttpRequestMethod
// -------------------------------------------------------------------------------------------------

func (r HttpRequestMethod) ToString() string {
	switch r {
	case METHOD_GET:
		return "GET"
	case METHOD_POST:
		return "POST"
	case METHOD_DELETE:
		return "DELETE"
	case METHOD_PATCH:
		return "PATCH"
	case METHOD_PUT:
		return "PUT"
	case METHOD_HEAD:
		return "HEAD"
	case METHOD_OPTIONS:
		return "OPTIONS"
	}
	return "UNKNOWN REQUEST METHOD"
}
