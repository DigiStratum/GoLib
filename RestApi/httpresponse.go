package restapi

// HTTP Response public interface
type HttpResponseIfc interface {
	GetBody() *string
	SetBody(body *string)
	GetStatus() HttpStatus
	SetStatus(status HttpStatus)
	GetHeaders() HttpHeadersIfc
}

type httpResponse struct {
	status		HttpStatus
	headers		HttpHeadersIfc
	body		*string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpResponse() HttpResponseIfc {
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

