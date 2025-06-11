package http

/*
Data set for HTTP request body data which supports multiple values for a given named property

This is typically form-encoded data, parsing required; note that this structure is important both
for our code to create requets and send to a server... and also for us to receive requests from
some client and act as the server. In the one case we must build the request and format it to send,
and in the other we must parse the request body data to extract the name-value pairs so that we
can respond in kind. In both cases, it is still a request body whether a request we are making or
one that we are receiving.

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
