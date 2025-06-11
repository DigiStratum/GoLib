package http

/*
Data set for HTTP request body data which supports multiple values for a given named property

This is typically form-encoded data, parsing required

TODO:
 * Implement based on the Headers/Builder; this differs from Headers in that there are no name-
   value pair standards, and no hard limits on sizes, etc

*/

type HttpRequestBodyDataIfc interface {
}

type httpRequestBodyData struct {
	// HTTP Request body data form-encoded supports name/-multivalue pairs
	// (i.e. a name can havev multiple, ordered values)
	bodyData map[string][]string
}
