package restapi

// HTTP Response public interface
type HttpResponseIfc interface {
	GetBody() *string
	SetBody(body *string)
	GetStatus() HttpStatus
	SetStatus(status HttpStatus)
	GetHeaders() *HttpHeaders
}

type HttpResponse struct {
	status		HttpStatus
	headers		*HttpHeaders
	body		*string
}

// Make a new one of these!
func NewHttpResponse() *httpResponse {
	hdrs := make(HttpHeaders)
	return &httpResponse{
		headers: &hdrs,
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

func (hr *httpResponse) GetHeaders() *HttpHeaders {
	return hr.headers
}

