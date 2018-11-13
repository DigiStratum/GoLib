package module

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
type controllerEPVMap		map[EndpointVersion]interface{}
type controllerEPPVMap		map[EndpointPattern]controllerEPVMap
type controllerEPMPVMap		map[EndpointMethod]controllerEPPVMap
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

var controller *Controller

func init() {
	controller = NewController()
}

// Get the singleton instance
func GetController() *Controller {
	return controller
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
	for _, endpoint := range (*ctrlr).endpoints {
		if ep, ok := endpoint.(EndpointIfc); ok {
			ep.Configure(endpoint, *serverConfig, *moduleConfig, *extraConfig)
		} else {
			// wot? Not an Endpoint!
			l.Error(fmt.Sprintf(
				"Controller{%s}.Configure(): Non-Endpoint given to Controller",
				moduleName,
			))
		}
	}
}

// Add an Endpoint to this Controller
// Endpoint uses this to self-register Endpoints at init() time
func (ctrlr *Controller) AddEndpoint(concreteEndpoint interface{}, name string, version string) {

	// Initialize the Endpoint
	ctrlr.endpoints = append(ctrlr.endpoints, concreteEndpoint)
	endpoint := concreteEndpoint.(EndpointIfc)
	endpoint.Init(concreteEndpoint, name, version)

	// Get the Endpoint's Pattern
	epp := EndpointPattern(endpoint.GetPattern())

	// See that the registry has an entry for each method/pattern/version for this Endpoint
	methods := endpoint.GetMethods()
	for _, method := range methods {
		// If this method isn't registered for this Controller, add it now
		epm := EndpointMethod(method)
		if _, ok := (*ctrlr.endpointMap)[epm]; !ok {
			(*ctrlr.endpointMap)[epm] = make(controllerEPPVMap)
		}

		// If this pattern isn't registered yet, add it now
		if _, ok := (*ctrlr.endpointMap)[epm][epp]; !ok {
			(*ctrlr.endpointMap)[epm][epp] = make(controllerEPVMap)
		}

		// If this version isn't registered yet, add it now
		if _, ok := (*ctrlr.endpointMap)[epm][epp][EndpointVersion(version)]; !ok {
			(*ctrlr.endpointMap)[epm][epp][EndpointVersion(version)] = endpoint
		}
	}
}

// Do any request pre-processing needed...
func (ctrlr *Controller) HandleRequest(request *rest.HttpRequest) *rest.HttpResponse {
	hlpr := rest.GetHelper()

	// Is the request method in our Endpoint registry?
	epm := EndpointMethod(request.GetMethod())
	if _, ok := (*ctrlr.endpointMap)[epm]; !ok {
		return hlpr.ResponseError(rest.STATUS_METHOD_NOT_ALLOWED)
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
	// TODO: Test more specific patterns before more general ones
	epm := EndpointMethod(request.GetMethod())
	for pattern, versions := range (*ctrlr.endpointMap)[epm] {
		// Test the pattern
		// ref: https://golang.org/pkg/regexp/#example_MatchString
		//l.Trace(fmt.Sprintf("Controller: Checking Pattern: '%s'", pattern))
		matches, err := ctrlr.matchesURI(string(pattern), relativeURI)
		if nil != err {
			l.Error(err.Error())
			return hlpr.ResponseError(rest.STATUS_INTERNAL_SERVER_ERROR)
		}
		if ! matches { continue }

		// Dispatch this request!
		for _, endpoint := range versions {
			// TODO: scan versions for version specified in X-Version header
			// Default to first listed version
			return endpointHandleRequest(endpoint, request)
		}
		return nil // UNHANDLED BY US
	}
	// If we fell through without finding a handler (Endpoint), then we're done for!
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
func (ctrlr *Controller) matchesURI(pattern string, URI string) (bool, error) {
	// Find the Regexp in the pattern cache
	var rxp	*regexp.Regexp
	var ok bool
	if rxp, ok = (*ctrlr.patternCache)[pattern]; !ok {
		// No!? Well then... Compile it and ADD it to the cache!
		var err error
		rxp, err = regexp.Compile(pattern)
		if nil != err {
			return false, errors.New(fmt.Sprintf(
				"Controller: Unexpected error compiling Regex pattern '%s'",
				pattern,
			))
		}
		(*ctrlr.patternCache)[pattern] = rxp
	}
	return rxp.MatchString(URI), nil
}

