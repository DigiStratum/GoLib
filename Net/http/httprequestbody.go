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
  * Make this iterable so that we can iterate over the name-value pairs
  * What is the purpose of the Size() function? It is attempting to figure the total size of the
    request body data as if it were formatted, however it is not implementing form encoding so
	it cannot be accurate. But do we need an accurate representation of size? Or do we need to
	know the size at all?

*/

type HttpRequestBodyIfc interface {
	Has(name string) bool
	GetNames() *[]string
	IsEmpty() bool
	Get(name string) *[]string
	Size() int
}

// HTTP Request body data form-encoded supports name/-multivalue pairs
// (i.e. a name can have multiple, ordered values)
type httpRequestBodyData map[string][]string

type httpRequestBody struct {
	body httpRequestBodyData
}

// -------------------------------------------------------------------------------------------------
// HttpRequestBodyIfc
// -------------------------------------------------------------------------------------------------

// DO we have the named property?
func (r *httpRequestBody) Has(name string) bool {
	_, ok := r.body[name]
	return ok
}

// Get the complete set of property names
func (r *httpRequestBody) GetNames() *[]string {
	names := make([]string, 0)
	for name, _ := range r.body {
		names = append(names, name)
	}
	return &names
}

// Are there NO body properties set?
func (r *httpRequestBody) IsEmpty() bool {
	return len(r.body) == 0
}

func (r *httpRequestBody) Get(name string) *[]string {
	if values, ok := r.body[name]; ok {
		return &values
	}
	return nil
}

/*
A set of body data will render out as a text block like:
---
name1=value1&name=value2&name=value3name2=value1&value2
---
But we consolidate duplicate names into a single name, multi-value memory structure, so we need to
re-expand into the original form for purposes of assessing the overall length. Note that the
leading ampersand is optinonal, but we will add the extra byte count for it.

Our Size() function will return the length of this text block, including the separators.
*/
func (r *httpRequestBody) Size() int {
	l := 0
	for name, values := range r.body {
		for _, value := range values {
			l += len(name) + 1  // Name length plus '&' pair separator
			l += len(value) + 1 // Value length plus '=' name-value separator
		}
	}
	return l
}
