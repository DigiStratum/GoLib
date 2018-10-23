package module

import(
	"fmt"
	"regexp"
	"errors"

	lib "../../golib"
	rest "../restapi"
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

type AbstractControllerIfc interface {
	Configure(serverConfig lib.Config, moduleConfig lib.Config)
	AddEndpoint(endpoint AbstractEndpointIfc)
	HandleRequest(request rest.HttpRequest) rest.HttpResponse
}

type AbstractController struct {
	securityPolicy		SecurityPolicy	// Module-wide SecurityPolicy
	serverConfig		lib.Config		// Server configuration cache
	moduleConfig		lib.Config		// Module configuration cache
	endpoints		controllerEPMPVMap	// Registry of all our Endpoints
	patternCache	regexpCache		// Compiled Regex Endpoint pattern cache
}

// Make a new one!
func NewController() *AbstractController {
	ctrlr := AbstractController{
		serverConfig:		lib.NewConfig(),
		moduleConfig:		lib.NewConfig(),
		endpoints:		make(controllerEPMPVMap),
		patternCache:	make(regexpCache),
	}
	return &ctrlr
}

// Module passes its own SecurityPolicy to us for reference
func (ctrlr *AbstractController) SetSecurityPolicy(securityPolicy SecurityPolicy) {
	ctrlr.securityPolicy = securityPolicy
}

// Module initializes a Controller
func (ctrlr *AbstractController) Configure(serverConfig lib.Config, moduleConfig lib.Config) {
	l := lib.GetLogger()
	l.Trace("Controller: Configure")

	// Capture the server and module configuration data for future reference
	ctrlr.serverConfig = serverConfig
	ctrlr.moduleConfig = moduleConfig

	// Configure the endpoints
	for _, patterns := range ctrlr.endpoints {
		for _, versions := range patterns {
			for _, endpoint := range versions {
				if ep, ok := endpoint.(AbstractEndpointIfc); ok {
					ep.Configure(serverConfig, moduleConfig)
				} else {
					// wot? Not an AbstractEndpoint!
					l.Error(fmt.Sprintf(
						"Controller: Non-Endpoint given to Controller in Module '%s'",
						ctrlr.moduleConfig.Get("module.name"),
					))
				}
			}
		}
	}
}

// Add a single Endpoint to this Controller
func (ctrlr *AbstractController) AddEndpoint(endpoint AbstractEndpointIfc) {
	// Get the Endpoint's Pattern/Version
	epp := EndpointPattern(endpoint.GetPattern())
	epv := EndpointVersion(endpoint.GetVersion())

	// See that the registry has an entry for each method/pattern/version for this Endpoint
	methods := endpoint.GetMethods()
	for _, method := range methods {
		// If this method isn't registered yet, add it now
		epm := EndpointMethod(method)
		if _, ok := ctrlr.endpoints[epm]; !ok {
			ctrlr.endpoints[epm] = make(controllerEPPVMap)
		}

		// If this pattern isn't registered yet, add it now
		if _, ok := ctrlr.endpoints[epm][epp]; !ok {
			ctrlr.endpoints[epm][epp] = make(controllerEPVMap)
		}

		// If this version isn't registered yet, add it now
		if _, ok := ctrlr.endpoints[epm][epp][epv]; !ok {
			ctrlr.endpoints[epm][epp][epv] = endpoint
		}
	}
}

// Do any request pre-processing needed...
func (ctrlr *AbstractController) HandleRequest(request rest.HttpRequest) *rest.HttpResponse {
	hlpr := rest.GetHelper()

	// Is the request method in our Endpoint registry?
	epm := EndpointMethod(request.GetMethod())
	if _, ok := ctrlr.endpoints[epm]; !ok {
		return hlpr.ResponseError(rest.STATUS_METHOD_NOT_ALLOWED)
	}

	// Will our Module SecurityPolicy reject this Request?
	if rej := ctrlr.securityPolicy.HandleRejection(request); nil != rej { return rej } // REJECT!

	return ctrlr.dispatchRequest(&request)
}

// Dispatch the request to an Endpoint
func (ctrlr *AbstractController) dispatchRequest(request *rest.HttpRequest) *rest.HttpResponse {
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
	for pattern, versions := range ctrlr.endpoints[epm] {
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
	if ep, ok := endpoint.(AbstractEndpointIfc); ok {
		l.Trace(fmt.Sprintf(
			"[%s] Controller: Selected Endpoint: '%s'",
			ctx.GetRequestId(),
			ep.GetName(),
		))
		return ep.HandleRequest(*request, &ep)
	}
	l.Error(fmt.Sprintf(
		"[%s] Controller: Unexpected error converting to Abstract Endpoint",
		ctx.GetRequestId(),
	))
	hlpr := rest.GetHelper()
	return hlpr.ResponseError(rest.STATUS_INTERNAL_SERVER_ERROR)
}

// Use a pattern cache of compiled RegExp's to match the URI
func (ctrlr *AbstractController) matchesURI(pattern string, URI string) (bool, error) {
	// Find the Regexp in the pattern cache
	var rxp	*regexp.Regexp
	var ok bool
	if rxp, ok = ctrlr.patternCache[pattern]; !ok {
		// No!? Well then... Compile it and ADD it to the cache!
		var err error
		rxp, err = regexp.Compile(pattern)
		if nil != err {
			return false, errors.New(fmt.Sprintf(
				"Controller: Unexpected error compiling Regex pattern '%s'",
				pattern,
			))
		}
		ctrlr.patternCache[pattern] = rxp
	}
	return rxp.MatchString(URI), nil
}

