package http

/*

Configurable, reusable client for making HTTP Requests

TODO:
 * Can we use Go Routines to set up an async pool of sorts so that we can have multiple requests in
   flight?
*/

import(
	"fmt"
	"net/http"
)

type HttpClientIfc interface {
	GetRequestResponse(request HttpRequestIfc) (*httpResponse, error)
}

type httpClient struct {
}

func NewHttpClient() *httpClient {
	r := httpClient{}
	return &r
}

func (r *httpClient) GetRequestResponse(request HttpRequestIfc) (*httpResponse, error) {
	// FIXME: Take the request details, populate net/http, fire off the request, get the response, and transform it back to httpResponse.
	// TODO: Also accept any sort of config details from receiver such as timeout settings.
	switch (request.GetMethod()) {
		// Return an internal error instead of 400 response as if from server to clarify error origin
		case METHOD_GET:
			// FIXME - bring in all the other juicy details from the request structure as well
			// TODO: Capture the response
			http.Get(request.GetURL())
	}
	return nil, fmt.Errorf("Unknown Request Method")
}

