package http

/*

Configurable, reusable client for making HTTP Requests

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
}

func NewHttpClient() *httpClient {
	r := httpClient{}
	return &r
}

func (r *httpClient) GetRequestResponse(request HttpRequestIfc) (*httpResponse, error) {
	// FIXME: Take the request details, populate net/http, fire off the request, get the response, and transform it back to httpResponse.

	// TODO: Accept any sort of config details from receiver such as timeout settings and override defaults

	// TODO: Transport & Client should be created once and reused; this means that either
	// HttpClient itself should be a singleton, or at least have an extended lifecycle with
	// these bits configured during the initialization phase.

	// ref: https://pkg.go.dev/net/http
	transport := &gohttp.Transport{
		//MaxIdleConns:       10,
		//IdleConnTimeout:    30 * time.Second,
		//DisableCompression: true,
		IdleConnTimeout:    CLIENT_DEFAULT_TIMEOUT_IDLE * time.Second,
		DisableCompression: CLIENT_DEFAULT_COMPRESSION_DISABLE,
	}


	client := &gohttp.Client{
		//CheckRedirect: redirectPolicyFunc,
		Transport: transport,
	}
	hlpr := GetHelper()
	clientReq, err := http.NewRequest(
		hlpr.GetHttpRequestMethodText(request.GetMethod()),
		request.GetURL(),
		nil,
	)

	// Set up request headers
	headers := request.GetHeaders()
	headerNames := headers.GetHeaderNames()
	for _, headerName := range *headerNames {
		//req.Header.Add("If-None-Match", `W/"wyzzy"`)
		req.Header.Add(headerName, headers.Get(headerName))
	}

	resp, err := client.Do(req)
	/*
	switch (request.GetMethod()) {
		// Return an internal error instead of 400 response as if from server to clarify error origin
		case METHOD_GET:
			// FIXME - bring in all the other juicy details from the request structure as well
			// TODO: Capture the response
			result, err := http.Get(request.GetURL())
	}

	*/
	return nil, fmt.Errorf("Unknown Request Method")
}

