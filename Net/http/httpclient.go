package http

/*

Configurable, reusable client for making HTTP Requests

REF:
 * https://pkg.go.dev/net/http

 TODO:
 * Can we use Go Routines to set up an async pool of sorts so that we can have multiple requests in
   flight?
*/

import(
	"fmt"
	gohttp "net/http"
)

const CLIENT_DEFAULT_TIMEOUT_IDLE		30
const CLIENT_DEFAULT_COMPRESSION_DISABLE	false

type HttpClientIfc interface {
	GetRequestResponse(request HttpRequestIfc) (*httpResponse, error)
}

type httpClient struct {
	client			*gohttp.Client
}

func NewHttpClient() *httpClient {
	r := httpClient{
		client:		&gohttp.Client{
			transport:	&gohttp.Transport{
				IdleConnTimeout:    CLIENT_DEFAULT_TIMEOUT_IDLE * time.Second,
				DisableCompression: CLIENT_DEFAULT_COMPRESSION_DISABLE,
			},
		},
	}
	return &r
}

// Initialize a request/client, fire it off, get the response, and transform it back to httpResponse
func (r *httpClient) GetRequestResponse(request HttpRequestIfc) (*httpResponse, error) {

	// Transform the Request structure
	request, err := r.fromHttpRequest(httpRequest)
	if nil != err { return nil, err }

	// Make the request and capture the response
	response, err := client.Do(req)
	if nil != err { return nil, err }

	// Transform the Response structure
	httpResponse, err := r.toHttpResponse(response)

	returun httpResponse, err
}

// Data transform from our own HttpRequestIfc to Go net/http::Request
func (r *httpClient) fromHttpRequest(httpRequest HttpRequestIfc) (*gohttp.Request, error) {
	hlpr := GetHelper()
	request, err := http.NewRequest(
		hlpr.GetHttpRequestMethodText(request.GetMethod()),
		request.GetURL(),
		nil,
	)

	// Set up request headers
	headers := request.GetHeaders()
	headerNames := headers.GetHeaderNames()
	for _, headerName := range *headerNames {
		request.Header.Add(headerName, headers.Get(headerName))
	}

	return &request, nil
}

// Data transform to our own HttpResponseIfc from Go net/http::Response
func (r *httpClient) toHttpResponse(response gohttp.Response) (*httpResponse, error) {
	// FIXME: transform resp into an HttpResponseIfc
	return nil, nil
}


