package module

import (
	rest "github.com/DigiStratum/GoLib/RestApi"
)

type SecurityPolicy struct {
	requireAuthentication	bool
}

func NewSecurityPolicy() *SecurityPolicy {
	sp := SecurityPolicy{
		requireAuthentication:	false,
	}
	return &sp
}

// Set whether authentication is required
func (sp *SecurityPolicy) SetRequireAuthentication(isRequired bool) {
	sp.requireAuthentication = isRequired
}

// Get whether authentication is required
func (sp *SecurityPolicy) RequiresAuthentication() bool {
	return sp.requireAuthentication
}

// How to handle rejection: Get an appropriate HttpResponse for any SecurityPolicy rejection; nil if ok
func (sp SecurityPolicy) HandleRejection(request *rest.HttpRequest) *rest.HttpResponse {
	hlpr := rest.GetHelper()

	// If authentication is required...
	if sp.RequiresAuthentication() {
		// ... make sure the request has provided an Authorization header
		requestHeaders := request.GetHeaders()
		authHeader := requestHeaders.Get("authorization")
		if "" == authHeader {
			return hlpr.ResponseError(rest.STATUS_UNAUTHORIZED)
		}
	}

	return nil
}

