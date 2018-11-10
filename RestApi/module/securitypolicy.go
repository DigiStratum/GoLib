package module

/*

SecurityPolicy - enforce HTTP request/response secuirty based on configured constraints

TODO: Expand to actually check HttpRequest auth with configured authenticator(s)

*/

import (
	rest "github.com/DigiStratum/GoLib/RestApi"
)

type SecurityPolicy struct {
	requireAuthentication	bool
}

// Make a new one of these!
func NewSecurityPolicy(config *lib.Config) *SecurityPolicy {

	// By default we do nothing
	sp := SecurityPolicy{
		requireAuthentication:	false,
	}

	// By configuration, we start enabling things...
	if len(config) > 0 {
		if "true" == config.Get("isrequired") {
			sp.SetRequireAuthentication(true)
		}
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

