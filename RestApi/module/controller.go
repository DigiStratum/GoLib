package module

/*

The Controller is where Module-specific HTTP Request handling takes place.

Endpoints are mapped with regular expression patterns. Multiple regular expression patterns could
potentially match the same string, in this case the request URI. For example a pattern matching ".*"
could more generally match the same thing as "specific_thing/\d+". To counter this we use a scoring
system that subtracts all regex-matched segmements from the complete URI match - if all patterns
match the entire URI, the more specific matches will have fewer characters in the varigable segments
,atched by the regex.

More thoughts are here:
ref: https://cs.stackexchange.com/questions/10786/how-to-find-specificity-of-a-regex-match

TODO: Do we really need/want version support for Endpoint?

*/

import(
	"fmt"
	"path"
	"strings"

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
type endpointContainer struct {
	endpointMPV	EndpointIfc			// The endpoint itself
	sequence	int				// Sequentialize mappings to ensure pattern matching sequencing
}

type controllerEPVMap	map[string]endpointContainer	// Endpoint Version map
type controllerEPPVMap	map[string]controllerEPVMap	// Endpoint Pattern map
type controllerEPMPVMap	map[string]controllerEPPVMap	// Endpoint Method map

type Controller struct {
	securityPolicy	*SecurityPolicy		// Module-wide SecurityPolicy
	serverConfig	*lib.Config		// Server configuration cache
	moduleConfig	*lib.Config		// Module configuration cache
	extraConfig	*lib.Config		// Extra configuration for Endpoints
	endpointMap	*controllerEPMPVMap	// Map of all our Endpoints
	endpoints	[]interface{}		// Collection of distinct concrete endpoints
}

// Make a new one of these!
func NewController() *Controller {
	c := make(controllerEPMPVMap)
	return &Controller{
		serverConfig:	lib.NewConfig(),
		moduleConfig:	lib.NewConfig(),
		extraConfig:	lib.NewConfig(),
		endpointMap:	&c,
		endpoints:	make([]interface{}, 0),
	}
}

// Module passes its own SecurityPolicy to us for reference
func (ctrlr *Controller) SetSecurityPolicy(securityPolicy *SecurityPolicy) {
	ctrlr.securityPolicy = securityPolicy
}

// Wrap log messages with Controller context to reduce boilerplate elsewhere
func (ctrlr *Controller) wrapLog(msg string) string {
	return fmt.Sprintf("Controller{%s}.%s", ctrlr.moduleConfig.Get("name"), msg)
}

// Module initializes a Controller
func (ctrlr *Controller) Configure(serverConfig *lib.Config, moduleConfig *lib.Config, extraConfig *lib.Config) {
	l := lib.GetLogger()
	l.Trace(ctrlr.wrapLog("Configure()"))

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
			l.Error(ctrlr.wrapLog("Configure(): Non-Endpoint given to Controller"))
		}
	}
}

// See that the map has an entry for each method/pattern/version for this Endpoint
func (ctrlr *Controller) mapEndpoint(endpoint EndpointIfc) {
	// Get the Endpoint's Pattern; we force it to match entire URI following Module prefix
	pattern := endpoint.GetPattern()
	methods := endpoint.GetMethods()
	version := endpoint.GetVersion()
	for _, method := range methods {
		// If this method isn't registered for this Controller, add it now
		if _, ok := (*ctrlr.endpointMap)[method]; !ok {
			(*ctrlr.endpointMap)[method] = make(controllerEPPVMap)
		}

		// If this pattern isn't registered yet, add it now
		if _, ok := (*ctrlr.endpointMap)[method][pattern]; !ok {
			(*ctrlr.endpointMap)[method][pattern] = make(controllerEPVMap)
		}

		// If this version isn't registered yet, add it now
		if _, ok := (*ctrlr.endpointMap)[method][pattern][version]; !ok {
			(*ctrlr.endpointMap)[method][pattern][version] = endpointContainer{
				endpointMPV: endpoint,
				sequence: len(*ctrlr.endpointMap),
			}
		}
	}
}

// Do any request pre-processing needed...
func (ctrlr *Controller) HandleRequest(request *rest.HttpRequest) *rest.HttpResponse {
	lib.GetLogger().Trace(ctrlr.wrapLog(fmt.Sprintf(
		"[%s]HandleRequest(): %s %s",
		request.GetContext().GetRequestId(),
		request.GetMethod(),
		request.GetURL(),
	)))

	// Is the request method in our Endpoint registry?
	if _, ok := (*ctrlr.endpointMap)[request.GetMethod()]; !ok {
		return rest.GetHelper().ResponseError(rest.STATUS_METHOD_NOT_ALLOWED)
	}

	// Will our Module SecurityPolicy reject this Request?
	if rej := ctrlr.securityPolicy.HandleRejection(request); nil != rej { return rej }

	response := ctrlr.dispatchRequest(request)

	// If the response was 404, see if an alternate 
	if response.GetStatus() == rest.STATUS_NOT_FOUND {
		altResponse := ctrlr.alternateRequest(request)
		if nil != altResponse { return altResponse }
	}

	// Otherwise, use the original response
	return response
}

// Retry requests with missing trailing slash
func (ctrlr *Controller) alternateRequest(request *rest.HttpRequest) *rest.HttpResponse {
	// If the last component of the request URI is empty, then it already trails a '/'
	_, fileName := path.Split(request.GetURI())
	if len(fileName) == 0 { return nil }

	// If the last component of the request URI has a file extension then assume not a dir
	pos := strings.LastIndex(fileName, ".")
	if 0 <= pos { return nil }

	// Try again witih a trailing '/'
	request.SetURI(fmt.Sprintf("%s/", request.GetURI()))
	request.SetURL(fmt.Sprintf("%s/", request.GetURL()))
	response := ctrlr.dispatchRequest(request)

	// If the response is still an error, then things not improving so keep original response
	if response.GetStatus() >= rest.STATUS_BAD_REQUEST { return nil }

	// If the request method was idempotent then redirect, otherwise return trhe response we got
	if request.IsIdempotentMethod() {
		// Redirect to the new location so that everything will
		// relativize correctly. This happens when the user types
		// In a directory, but doesn't put the trailing slash on it.
		hlpr := rest.GetHelper()
		return hlpr.ResponseRedirect(request.GetURL())
	}
	return response
}

// Dispatch the request to an Endpoint
func (ctrlr *Controller) dispatchRequest(request *rest.HttpRequest) *rest.HttpResponse {
	hlpr := rest.GetHelper()
	l := lib.GetLogger()

	// Strip the server/module components off the beginning of the URI
	ctx := request.GetContext()
	requestURI := request.GetURI()
	l.Trace(ctrlr.wrapLog(fmt.Sprintf(
		"[%s]dispatchRequest() - Dispatching: '%s'",
		ctx.GetRequestId(),
		requestURI,
	)))

	endpoint := ctrlr.findBestMatchingEndpointForURI(request)
	if nil != endpoint {
		//l.Trace(fmt.Sprintf("\tpassing to endpointHandleRequest for sequence %d", endpoint.sequence))
		return ctrlr.endpointHandleRequest((*endpoint).endpointMPV, request)
	}

	// We fell through without a match; last chance: if the REQUEST URI is root of the server/module...
	if requestURI == ctx.GetPrefixPath() {
		// ... and if the full REQUEST URL does not end with '/'...
		requestURL := request.GetURL();
		if ! strings.HasSuffix(requestURL, "/") {
			// ... it's a non-conforming REQUEST URI; try again, but with trailing '/'
			return hlpr.ResponseRedirectPermanent(requestURL + "/")
		}
	}

	// If we fell through without finding a handler then we're done for!
	return hlpr.ResponseError(rest.STATUS_NOT_FOUND)
}

// Find which Endpoint's pattern matches this request URI - the BEST match, not just ANY match
func (ctrlr *Controller) findBestMatchingEndpointForURI(request *rest.HttpRequest) *endpointContainer {
	bestScore := 0
	var bestEndpoint *endpointContainer
	for _, versions := range (*ctrlr.endpointMap)[request.GetMethod()] {
		for version, endpoint := range versions {
			// TODO: Add support for client to specify a version with X-Version header? Maybe we don't need versions at all?
			matches := endpoint.endpointMPV.GetRequestMatches(request)
			if (nil == matches) || (len(matches) == 0) {
				lib.GetLogger().Crazy(ctrlr.wrapLog(fmt.Sprintf(
					"findBestMatchingEndpointForURI() - Pattern: '%s' for version '%s' - No Matches!",
					endpoint.endpointMPV.GetPattern(),
					version,
				)))
				continue
			}

			// Calculate a score for this pattern to determine how well it matches
			penalty := 0
			points := 0
			// When scoring for specificity, static matches score points and variability is penalized
			// May the best pattern win! Score = len(request.URI) - len(all matches)
			for index, match := range matches {
				if 0 == index {
					// The first match in matches is the completely matched string
					points = len(match)
				} else {
					// All other matches count against the score since they represent variability
					penalty += len(match)
				}
			}
			score := points - penalty

			// If current pattern scores better than best thus far
			if score >= bestScore {
				bestScore = score
				bestEndpoint = &endpoint
			}

			lib.GetLogger().Crazy(ctrlr.wrapLog(fmt.Sprintf(
				"findBestMatchingEndpointForURI() - Pattern: '%s' for version '%s', (matches: %d, points: %d, penalty: %d, score: %d)",
				endpoint.endpointMPV.GetPattern(),
				version,
				len(matches),
				points,
				penalty,
				score,
			)))
		}
	}

	// Note: Endpoint can call its own matches method to geth named path parameters 
	return bestEndpoint
}

// Wrap HTTP Request to send to Endpoint for handling
func (ctrlr *Controller) endpointHandleRequest(endpoint EndpointIfc, request *rest.HttpRequest) *rest.HttpResponse {
	ctx := request.GetContext()
	lib.GetLogger().Trace(ctrlr.wrapLog(fmt.Sprintf(
		"[%s]endpointHandleRequest() - Controller: Selected Endpoint: '%s'",
		ctx.GetRequestId(),
		endpoint.GetName(),
	)))
	res := endpoint.HandleRequest(request, endpoint)
	ctrlr.mergeDefaultResponseHeaders(res, ctx.GetRequestId())
	return res
}

// Merge default response headers into OK responses if not already supplied
// TODO: Make this more granular to be endpoint-specific (override)
func (ctrlr *Controller) mergeDefaultResponseHeaders(response *rest.HttpResponse, requestId string) {
	// Only OK (2xx) responses want default headers, else nothing to do
	if rest.GetHelper().IsStatus2xx(response.GetStatus()) { return }

	// Define default response header set
	mergeHeaders := ctrlr.getDefaultResponseHeaders()

	// Override default endpoints with module-specific configured response headers
	moduleHeaders := ctrlr.moduleConfig.GetSubset("headers.")

	l := lib.GetLogger()
	// For each of the default response headers...
	for kvp := range moduleHeaders.IterateChannel() {
		l.Crazy(ctrlr.wrapLog(fmt.Sprintf(
			"[%s]mergeDefaultResponseHeaders() - Controller: Override Response Header: '%s' = '%s'",
			requestId,
			kvp.Key,
			kvp.Value,
		)))
		(*mergeHeaders)[kvp.Key] = kvp.Value
	}

	// Snag the response headers to check and modify
	headers := response.GetHeaders()

	// For each of the default response headers...
	for name, value := range *mergeHeaders {
		// If this header is not already in the set (i.e. endpoint's version wins)...
		if _, ok := (*headers)[name]; !ok {
			// Add it with the default value!
			(*headers)[name] = value
		}
	}

	// TODO: Add other response-dependent default headers here like e-tag, Last-Modified, etc.
}

// Get the default response header set
// TODO: Make cache-control expiration configurable
// ref: https://www.keycdn.com/blog/http-security-headers
// ref: https://www.keycdn.com/support/content-security-policy
// ref: https://www.keycdn.com/blog/http-cache-headers
// ref: https://content-security-policy.com/
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
func (ctrlr *Controller) getDefaultResponseHeaders() *map[string]string {
	headers := map[string]string{
		"cache-control": "max-age=43200,public,must-revalidate,no-transform",
		"x-frame-options": "SAMEORIGIN",
		"x-content-type-options": "nosniff",
		"x-xss-protection": "1; mode=block",
		"vary": "Accept-Encoding, Origin",
		"strict-transport-security": "max-age=31536000; includeSubdomains; preload",
		"access-control-allow-origin": "*",
		"content-security-policy": "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; frame-ancestors 'none';",
		"feature-policy": "autoplay 'none'; camera 'none'",
	}
	return &headers
}

