package http

/*
Data set for HTTP response body data which supports multiple values for a given named property

This is typically form-encoded data, parsing required

*/

//type HttpBodyData map[string][]string

type httpBodyData struct {
	bodyData map[string][]string
}
