package module

/*

The Controller is where Module-specific HTTP Request handling takes place.

TODO: Test more specific patterns before more general ones
ref: https://cs.stackexchange.com/questions/10786/how-to-find-specificity-of-a-regex-match
*/

import(
	"fmt"
	"regexp"
	"errors"

	lib "github.com/DigiStratum/GoLib"
	rest "github.com/DigiStratum/GoLib/RestApi"
)

// These are stored in this sequence for expedient dispatching of a request
// Because we use regular expressions to match the pattern of which Endpoint should handle a given
// request, it helps to reduce the number of regular expressions we try to find a match on by using
// the HTTP request method first to cut down the list, then trying only the patterns within the
// matching request method. From there, we will probably typically only have on version (and if we
// end up ALWAYS only being one version because we are awesome at ensuring non-breaking changes then
// we can remove versioning support at some point), so if there is only one, we will deliver it,
// even if the request specifies some version which may/may not match. 
type controllerEPVMap		map[string]interface{}		// Endpoint Version map
type controllerEPPVMap		map[string]controllerEPVMap	// Endpoint Pattern map
type controllerEPMPVMap		map[string]controllerEPPVMap	// Endpoint Method map
type regexpCache		map[string]*regexp.Regexp

type Controller struct {
	securityPolicy	*SecurityPolicy		// Module-wide SecurityPolicy
	serverConfig	*lib.Config		// Server configuration cache
	moduleConfig	*lib.Config		// Module configuration cache
	extraConfig	*lib.Config		// Extra configuration for Endpoints
	patternCache	*regexpCache		// Compiled Regex Endpoint pattern cache
	endpointMap	*controllerEPMPVMap	// Map of all our Endpoints
	endpoints	[]interface{}		// Collection of distinct concrete endpoints
}

// Make a new one of these!
func NewController() *Controller {
	c := make(controllerEPMPVMap)
	r := make(regexpCache)
	return &Controller{
		serverConfig:	lib.NewConfig(),
		moduleConfig:	lib.NewConfig(),
		extraConfig:	lib.NewConfig(),
		patternCache:	&r,
		endpointMap:	&c,
		endpoints:	make([]interface{}, 0),
	}
}

// Module passes its own SecurityPolicy to us for reference
func (ctrlr *Controller) SetSecurityPolicy(securityPolicy *SecurityPolicy) {
	ctrlr.securityPolicy = securityPolicy
}

// Module initializes a Controller
func (ctrlr *Controller) Configure(serverConfig *lib.Config, moduleConfig *lib.Config, extraConfig *lib.Config) {
	moduleName := ctrlr.moduleConfig.Get("name")
	l := lib.GetLogger()
	l.Trace(fmt.Sprintf("Controller{%s}.Configure()", moduleName))

	// Capture the server and module configuration data for future reference
	ctrlr.serverConfig = serverConfig
	ctrlr.moduleConfig = moduleConfig
	ctrlr.extraConfig = extraConfig

	// Configure all the Endpoints
	endpoints := GetRegistry().DrainEndpoints()
	ctrlr.endpoints = *endpoints
	for _, endpoint := range (*ctrlr).endpoints {
		if endpointIfc, ok := endpoint.(EndpointIfc); ok {
			endpointIfc.Init(endpoint)
			endpointIfc.Configure(endpoint, *serverConfig, *moduleConfig, *extraConfig)
			ctrlr.mapEndpoint(endpointIfc)
		} else {
			// wot? Not an Endpoint!
			l.Error(fmt.Sprintf(
				"Controller{%s}.Configure(): Non-Endpoint given to Controller",
				moduleName,
			))
		}
	}
}

// See that the map has an entry for each method/pattern/version for this Endpoint
func (ctrlr *Controller) mapEndpoint(endpoint EndpointIfc) {
	// Get the Endpoint's Pattern; we force it to match entire URI following Module prefix
	epp := endpoint.GetPattern()
	methods := endpoint.GetMethods()
	version := endpoint.GetVersion()
	for _, method := range methods {
		// If this method isn't registered for this Controller, add it now
		epm := method
		if _, ok := (*ctrlr.endpointMap)[epm]; !ok {
			(*ctrlr.endpointMap)[epm] = make(controllerEPPVMap)
		}

		// If this pattern isn't registered yet, add it now
		if _, ok := (*ctrlr.endpointMap)[epm][epp]; !ok {
			(*ctrlr.endpointMap)[epm][epp] = make(controllerEPVMap)
		}

		// If this version isn't registered yet, add it now
		if _, ok := (*ctrlr.endpointMap)[epm][epp][version]; !ok {
			(*ctrlr.endpointMap)[epm][epp][version] = endpoint
		}
	}
}

// Do any request pre-processing needed...
func (ctrlr *Controller) HandleRequest(request *rest.HttpRequest) *rest.HttpResponse {
	lib.GetLogger().Trace(fmt.Sprintf(
		"[%s] Controller{%s)}.HandleRequest(): %s %s",
		request.GetContext().GetRequestId(),
		ctrlr.moduleConfig.Get("name"),
		request.GetMethod(),
		request.GetURL(),
	))

	// Is the request method in our Endpoint registry?
	epm := request.GetMethod()
	if _, ok := (*ctrlr.endpointMap)[epm]; !ok {
		return rest.GetHelper().ResponseError(rest.STATUS_METHOD_NOT_ALLOWED)
	}

	// Will our Module SecurityPolicy reject this Request?
	if rej := ctrlr.securityPolicy.HandleRejection(request); nil != rej { return rej }

	return ctrlr.dispatchRequest(request)
}

// Dispatch the request to an Endpoint
func (ctrlr *Controller) dispatchRequest(request *rest.HttpRequest) *rest.HttpResponse {
	hlpr := rest.GetHelper()
	l := lib.GetLogger()
	// Strip the server/module components off the beginning of the URI
	ctx := request.GetContext()
	requestUri := request.GetURI()
	relativeURI := requestUri[len(ctx.GetPrefixPath()):]
	l.Trace(fmt.Sprintf(
		"[%s] Controller: Dispatching: '%s'",
		ctx.GetRequestId(),
		relativeURI,
	))

	// Find which Endpoint's pattern matches this request URI
	// Note: we will find the BEST match, not just any match
	epm := request.GetMethod()
	bestScore := 0
	var bestVersions controllerEPVMap
	for pattern, versions := range (*ctrlr.endpointMap)[epm] {
		l.Trace(fmt.Sprintf("Controller: Checking Pattern: '%s'", pattern))
		matches, err := ctrlr.getUriMatches(pattern, relativeURI)
		if nil != err {
			l.Error(err.Error())
			return hlpr.ResponseError(rest.STATUS_INTERNAL_SERVER_ERROR)
		}
		if (nil == matches) || (len(matches) == 0) { continue }

		// Calculate a score for this pattern to determine how well it matches
		score := 0
		for _, match := range matches { score += len(match) }

		// If the current pattern scores better than the best pattern thus far...
		if score > bestScore {
			// then make this pattern the new best pattern!
			bestScore = score
			bestVersions = versions
			// TODO: Capture the matches into the Request Context;
			// it will have Endpoint specific parametric breakdown!
		}
	}

	// Use the versions of the best pattern we've found, if any
	if nil != bestVersions {
		// Dispatch this request!
		for _, endpoint := range bestVersions {
			// TODO: scan versions for version specified in X-Version header
			// Default to first listed version
			return endpointHandleRequest(endpoint, request)
		}
		return nil // UNHANDLED BY US
	}

	// If we fell through without finding a handler then we're done for!
	return hlpr.ResponseError(rest.STATUS_NOT_FOUND)
}

// Pass this request to the supplied Endpoint
func endpointHandleRequest(endpoint interface{}, request *rest.HttpRequest) *rest.HttpResponse {
	ctx := request.GetContext()
	l := lib.GetLogger()
	if ep, ok := endpoint.(EndpointIfc); ok {
		l.Trace(fmt.Sprintf(
			"[%s] Controller: Selected Endpoint: '%s'",
			ctx.GetRequestId(),
			ep.GetName(),
		))
		return ep.HandleRequest(request, ep)
	}
	l.Error(fmt.Sprintf(
		"[%s] Controller: Unexpected error converting to  Endpoint",
		ctx.GetRequestId(),
	))
	hlpr := rest.GetHelper()
	return hlpr.ResponseError(rest.STATUS_INTERNAL_SERVER_ERROR)
}

// Use a pattern cache of compiled RegExp's to match the URI
func (ctrlr *Controller) getUriMatches(pattern string, URI string) ([]string, error) {
	// Find the Regexp in the pattern cache
	var rxp	*regexp.Regexp
	var ok bool
	if rxp, ok = (*ctrlr.patternCache)[pattern]; !ok {
		// No!? Well then... Compile it and ADD it to the cache!
		var err error
		rxp, err = regexp.Compile("^" + pattern + "$")
		if nil != err {
			return nil, errors.New(fmt.Sprintf(
				"Controller: Unexpected error compiling Regex pattern '%s'",
				pattern,
			))
		}
		(*ctrlr.patternCache)[pattern] = rxp
	}
	return rxp.FindStringSubmatch(URI), nil
}

