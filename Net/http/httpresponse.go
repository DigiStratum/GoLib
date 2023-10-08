package http

// HTTP Response public interface
type HttpResponseIfc interface {
	GetBody() *string
	SetBody(body *string)
	GetStatus() HttpStatus
	SetStatus(status HttpStatus)
	GetHeaders() HttpHeadersIfc
	SetHeaders(headers HttpHeadersIfc)
}

type httpResponse struct {
	status		HttpStatus
	// TODO:	add protocol version which is part of the server response, NOT in headers!
	headers		HttpHeadersIfc
	body		*string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpResponse() *httpResponse {
	return &httpResponse{
		headers: NewHttpHeaders(),
	}
}

// -------------------------------------------------------------------------------------------------
// HttpResponseIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *httpResponse) GetBody() *string {
	return r.body
}

func (r *httpResponse) SetBody(body *string) {
	r.body = body
}

func (r *httpResponse) GetStatus() HttpStatus {
	return r.status
}

func (r *httpResponse) SetStatus(status HttpStatus) {
	r.status = status
}

func (r *httpResponse) GetHeaders() HttpHeadersIfc {
	return r.headers
}

func (r *httpResponse) SetHeaders(headers HttpHeadersIfc) {
	r.headers = headers
}

