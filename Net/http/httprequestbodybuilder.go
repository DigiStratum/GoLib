package http

/*

Builder for HTTP Request Body

TODO:
  * Add a Factory Function to produce an HttpRequestBody from form-encoded POST data
  * Consider a Factory Function to inherit from an existing HttpRequestBody
*/

type HttpRequestBodyBuilderIfc interface {
	Set(name string, values ...string)
	Merge(requestBody HttpRequestBodyIfc)
	GetHttpRequestBody() HttpRequestBodyIfc
}

// Name/value pair header map for Request or Response
type httpRequestBodyBuilder struct {
	requestBody *httpRequestBody
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpRequestBodyBuilder() *httpRequestBodyBuilder {
	r := httpRequestBodyBuilder{
		requestBody: &httpRequestBody{
			body: make(httpRequestBodyData),
		},
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// HttpRequestBodyBuilderIfc
// -------------------------------------------------------------------------------------------------

// Single-name, multi-value support
func (r *httpRequestBodyBuilder) Set(name string, values ...string) *httpRequestBodyBuilder {
	// If the named header is not set, then create it

	if _, ok := r.requestBody.body[name]; !ok {
		r.requestBody.body[name] = make([]string, 0)
	}

	for _, value := range values {
		r.requestBody.body[name] = append(r.requestBody.body[name], value)
	}
	return r
}

func (r *httpRequestBodyBuilder) Merge(requestBody *httpRequestBody) *httpRequestBodyBuilder {
	if requestBody != nil {
		for name, values := range requestBody.body {
			// Use Set() to merge provided values with existing, instead of overwriting
			r.Set(name, values...)
		}
	}
	return r
}

func (r *httpRequestBodyBuilder) GetHttpRequestBody() *httpRequestBody {
	return r.requestBody
}
