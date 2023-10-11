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
	"io"
	"fmt"
	"time"
	gohttp "net/http"
)

const CLIENT_DEFAULT_TIMEOUT_IDLE =		30
const CLIENT_DEFAULT_COMPRESSION_DISABLE =	false

type HttpClientIfc interface {
	GetRequestResponse(request HttpRequestIfc) (*httpResponse, error)
}

type httpClient struct {
	client			*gohttp.Client
}

func NewHttpClient() *httpClient {
	r := httpClient{
		client:		&gohttp.Client{
			Transport:	&gohttp.Transport{
				IdleConnTimeout:    CLIENT_DEFAULT_TIMEOUT_IDLE * time.Second,
				DisableCompression: CLIENT_DEFAULT_COMPRESSION_DISABLE,
			},
		},
	}
	return &r
}

// Initialize a request/client, fire it off, get the response, and transform it back to httpResponse
func (r *httpClient) GetRequestResponse(httpRequest HttpRequestIfc) (*httpResponse, error) {

	// Transform the Request structure
	request, err := r.fromHttpRequest(httpRequest)
	if nil != err { return nil, err }

	// Make the request and capture the response
	response, err := r.client.Do(request)
	if nil != err { return nil, err }

	// Transform the Response structure
	httpResponse, err := r.toHttpResponse(response)

	return httpResponse, err
}

// Data transform from our own HttpRequestIfc to Go net/http::Request
func (r *httpClient) fromHttpRequest(httpRequest HttpRequestIfc) (*gohttp.Request, error) {
	hlpr := GetHelper()
	request, err := gohttp.NewRequest(
		hlpr.GetHttpRequestMethodText(httpRequest.GetMethod()),
		httpRequest.GetURL(),
		nil,
	)
	if nil != err { return nil, err }

	// Set up request headers
	headers := httpRequest.GetHeaders()
	headerNames := headers.GetHeaderNames()
	for _, headerName := range *headerNames {
		request.Header.Add(headerName, headers.Get(headerName))
	}

	return request, nil
}

// Data transform to our own HttpResponseIfc from Go net/http::Response
func (r *httpClient) toHttpResponse(response *gohttp.Response) (*httpResponse, error) {
	hlpr := GetHelper()

	httpResponse := NewHttpResponse()
	httpResponse.SetStatus(hlpr.GetHttpStatus(response.StatusCode))

	// Capture the protocol version from the server response
	httpResponse.SetProtocolVersion(
		fmt.Sprintf("%d.%d", response.ProtoMajor, response.ProtoMinor),
	)

	// Transform response headers
	httpResponseHeaders := NewHttpHeaders()
	for name, values := range response.Header {
		for _, value := range values {
			httpResponseHeaders.Add(name, value)
		}
	}
	httpResponse.SetHeaders(httpResponseHeaders)

	// TODO: Transform response body
	if response.ContentLength > 0 {
		readbuf := make([]byte, 65536)
		// TODO: We should protect ourselves here with a reasonable default max, and allow
		// consumer to override; don't just allow whatever response size to consume all
		// available memory!
		bodybuf := make([]byte, response.ContentLength)
		var err error
		var readlen int
		var pos int64
		for ; pos < response.ContentLength; pos += int64(readlen) {
			readlen, err = response.Body.Read(bodybuf)
			// EOF means we're done. readlen *should* never be
			// 0, but trap just in case to prevent perpetual loop
			if (io.EOF == err) || (0 == readlen) { break }
			if nil != err { return nil, err }
			copy(bodybuf[pos:], readbuf[0:readlen - 1])
		}
		httpResponse.SetBinBody(&bodybuf)
	}

	return httpResponse, nil
}


