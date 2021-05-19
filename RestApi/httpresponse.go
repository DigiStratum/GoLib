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

// Make a new one of these!
func NewHttpResponse() HttpResponseIfc {
	return &httpResponse{
		headers: NewHttpHeaders(),
	}
}

func (hr *httpResponse) GetBody() *string {
	return hr.body
}

func (hr *httpResponse) SetBody(body *string) {
	hr.body = body
}

func (hr *httpResponse) GetStatus() HttpStatus {
	return hr.status
}

func (hr *httpResponse) SetStatus(status HttpStatus) {
	hr.status = status
}

func (hr *httpResponse) GetHeaders() HttpHeadersIfc {
	return hr.headers
}

