package http

import (
	gohttp "net/http"
)

/*

Builder for HTTP Request Body

TODO:
  * Add a Factory Function to produce an HttpRequestBody from form-encoded POST data
  * Consider a Factory Function to inherit from an existing HttpRequestBody
*/

type HttpRequestBodyBuilderIfc interface {
	Set(name string, values ...string) *httpRequestBodyBuilder
	Merge(requestBody HttpRequestBodyIfc) *httpRequestBodyBuilder
	GetHttpRequestBody() *httpRequestBody
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

// NewHttpRequestBodyBuilderFromRequest creates a new HttpRequestBodyBuilder from an http.Request
// It automatically parses form data if needed and populates the body with PostForm values
func NewHttpRequestBodyBuilderFromRequest(request *gohttp.Request) *httpRequestBodyBuilder {
	if request == nil {
		return nil
	}

	// Parse form data if it hasn't been parsed yet
	err := request.ParseForm()
	if err != nil {
		// If parsing fails, return nil
		// TODO: Would be better to pass this error back
		return nil
	}

	builder := NewHttpRequestBodyBuilder()

	// Process all form values from the PostForm data
	for name, values := range request.PostForm {
		builder.Set(name, values...)
	}

	return builder
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

func (r *httpRequestBodyBuilder) Merge(requestBody HttpRequestBodyIfc) *httpRequestBodyBuilder {
	if requestBody != nil {
		names := requestBody.GetNames()
		if names == nil {
			return r // No names to merge
		}
		for _, name := range *names {
			values := requestBody.Get(name)
			if values == nil {
				continue // No values to merge for this name
			}
			// Use Set() to merge provided values with existing, instead of overwriting
			r.Set(name, *values...)
		}
	}
	return r
}

func (r *httpRequestBodyBuilder) GetHttpRequestBody() *httpRequestBody {
	return r.requestBody
}
