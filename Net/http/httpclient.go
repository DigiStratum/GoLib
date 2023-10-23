package http

/*

Configurable, reusable client for making HTTP Requests

REF:
 * https://pkg.go.dev/net/http

 TODO:
 * Can we use Go Routines to set up an async pool of sorts so that we can have multiple requests in
   flight?
 * Make Configurable to override defaults

*/

import(
	"io"
	"fmt"
	"time"
	gohttp "net/http"

	"github.com/DigiStratum/GoLib/Process/startable"
	cfg "github.com/DigiStratum/GoLib/Config"
)

type HttpClientIfc interface {
	// Embedded interface(s)
	cfg.ConfigurableIfc
	startable.StartableIfc

	// Our own interface
	GetRequestResponse(request HttpRequestIfc) (*httpResponse, error)
}

// Exported to support embedding
type HttpClient struct {
	// Embedded struct(s)
	*cfg.Configurable
	*startable.Startable

	// Our own properties
	client			*gohttp.Client
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpClient() *HttpClient {
	r := HttpClient{}

	// TODO: This int/bool-as-string madness caused by Config only supporting string values; add multi-type support!
	r.Configurable = cfg.NewConfigurable(
		cfg.NewConfigItem("maxBodyLenKb").SetDefault("10240"),
		cfg.NewConfigItem("idleTimeoutMSec").SetDefault("30000"),
		cfg.NewConfigItem("disableCompression").SetDefault("false"),
	)

	// Declare Starter funcs
	r.Startable = startable.NewStartable(
		r.Configurable,
	)

	return &r
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *HttpClient) Start() error {
	// Call super
	if err := r.Startable.Start(); nil != err { return err }

	// Set up our net/http client
	// ref: https://pkg.go.dev/net/http#Client
	r.client = &gohttp.Client{
		Transport:	&gohttp.Transport{
			// FIXME: Read these properties from ConfigItem values
			IdleConnTimeout:    CLIENT_DEFAULT_TIMEOUT_IDLE * time.Second,
			DisableCompression: CLIENT_DEFAULT_COMPRESSION_DISABLE,
		},
		// TODO: Support for CheckRedirect
		// TODO: Support for CookieJar
		Timeout: 0, // nonoseconds for request timeout; 0 = no timeout
	}

	return nil
}

// -------------------------------------------------------------------------------------------------
// HttpClientIfc
// -------------------------------------------------------------------------------------------------

// Initialize a request/client, fire it off, get the response, and transform it back to httpResponse
func (r *HttpClient) GetRequestResponse(httpRequest HttpRequestIfc) (*httpResponse, error) {

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
func (r *HttpClient) fromHttpRequest(httpRequest HttpRequestIfc) (*gohttp.Request, error) {
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
func (r *HttpClient) toHttpResponse(response *gohttp.Response) (*httpResponse, error) {
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

	// Transform response body
	if response.ContentLength > 0 {
		// Don't just allow any response size to consume all available memory!
		if response.ContentLength > r.maxbody {
			return nil, fmt.Errorf(
				"Response body length (%d) > max (%d)",
				response.ContentLength,
				r.maxbody,
			)
		}

		readbuf := make([]byte, 65536)
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


