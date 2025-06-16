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

import (
	"fmt"
	"io"
	gohttp "net/http"
	"strconv"
	"time"

	cfg "github.com/DigiStratum/GoLib/Config"
	"github.com/DigiStratum/GoLib/Process/startable"
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
	client *gohttp.Client

	// Config
	maxBodyLenKb       int
	requestTimeout     time.Duration
	idleTimeout        time.Duration
	disableCompression bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpClient() *HttpClient {
	r := HttpClient{}

	// TODO: This int/bool-as-string madness caused by Config only supporting string values; add multi-type support!
	r.Configurable = cfg.NewConfigurable(
		cfg.NewConfigItem("maxBodyLenKb").SetDefault("10240").CaptureWith(
			func(value string) error {
				v, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				r.maxBodyLenKb = v
				return nil
			},
		),
		cfg.NewConfigItem("requestTimeout").SetDefault("60s").CaptureWith(
			func(value string) error {
				var err error
				r.requestTimeout, err = time.ParseDuration(value)
				return err
			},
		),
		cfg.NewConfigItem("idleTimeout").SetDefault("30s").CaptureWith(
			func(value string) error {
				var err error
				r.idleTimeout, err = time.ParseDuration(value)
				return err
			},
		),
		cfg.NewConfigItem("disableCompression").SetDefault("false").CaptureWith(
			func(value string) error {
				r.disableCompression = (value == "true")
				return nil
			},
		),
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
	if err := r.Startable.Start(); err != nil {
		return err
	}

	// Set up our net/http client
	// ref: https://pkg.go.dev/net/http#Client
	r.client = &gohttp.Client{
		Transport: &gohttp.Transport{
			IdleConnTimeout:    r.idleTimeout,
			DisableCompression: r.disableCompression,
		},
		// TODO: Support for CheckRedirect
		// TODO: Support for CookieJar
		// nonoseconds for request timeout; 0 = no timeout
		Timeout: r.requestTimeout,
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
	if err != nil {
		return nil, err
	}

	// Make the request and capture the response
	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}

	// Transform the Response structure
	httpResponse, err := r.toHttpResponse(response)

	return httpResponse, err
}

// Data transform from our own HttpRequestIfc to Go net/http::Request
func (r *HttpClient) fromHttpRequest(httpRequest HttpRequestIfc) (*gohttp.Request, error) {
	request, err := gohttp.NewRequest(
		httpRequest.GetMethod().ToString(),
		httpRequest.GetURL(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Set up request headers
	headers := httpRequest.GetHeaders()
	headerNames := headers.GetNames()
	for _, headerName := range *headerNames {
		headerValues := headers.Get(headerName)
		if headerValues == nil {
			continue
		}
		request.Header.Add(headerName, (*headerValues)[0])
	}

	return request, nil
}

// Data transform to our own HttpResponseIfc from Go net/http::Response
func (r *HttpClient) toHttpResponse(response *gohttp.Response) (*httpResponse, error) {
	httpResponseBuilder := NewHttpResponseBuilder().
		SetStatus(HttpStatusFromCode(response.StatusCode))

	// Capture the protocol version from the server response
	httpResponseBuilder.SetProtocolVersion(
		fmt.Sprintf("%d.%d", response.ProtoMajor, response.ProtoMinor),
	)

	// Transform response headers
	httpResponseHeadersBuilder := NewHttpHeadersBuilder()
	for name, values := range response.Header {
		for _, value := range values {
			httpResponseHeadersBuilder.Set(name, value)
		}
	}
	httpResponseBuilder.SetHeaders(httpResponseHeadersBuilder.GetHttpHeaders())

	// Transform response body
	if response.ContentLength > 0 {
		// Don't just allow any response size to consume all available memory!
		if response.ContentLength > (int64(r.maxBodyLenKb) * 1024) {
			return nil, fmt.Errorf(
				"Response body length (%d) > max (%dKB)",
				response.ContentLength,
				r.maxBodyLenKb,
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
			if (err == io.EOF) || (readlen == 0) {
				break
			}
			if err != nil {
				return nil, err
			}
			copy(bodybuf[pos:], readbuf[0:readlen-1])
		}
		httpResponseBuilder.SetBinBody(&bodybuf)
	}

	return httpResponseBuilder.GetHttpResponse(), nil
}
