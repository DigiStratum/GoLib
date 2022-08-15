package webui

// An HTML Page is a flattened rendering of a Document with all magic tags resolved from the
// supplied Config. The Scheme shall be concatenated into a single block and injected into the
// top of the rendered Document.

// TODO: Deliver Scheme.stylesheet as a separate CSS file, separately requested (but referenced by Document).
// TODO: Use an external cache (memcache/elasticache?) to make rendered Documents available across
//	 instances using Configuration fingerprinting for differentiation

import(
	"fmt"

	log "github.com/DigiStratum/GoLib/Logger"
	cfg "github.com/DigiStratum/GoLib/Config"
)

const DEREFERENCE_MAX_LOOPS = 5

type HtmlPage struct {
	context			*cfg.Config	// Contextual/Configuration data available for injection into Document
	scheme			*Scheme		// Page layout/styling/structural templating
	document		*string		// Cached, final rendered document
}

// Make a new one of these; require a Scheme, but Context could be nil if there's nothing useful to inject
func NewHtmlPage(scheme *Scheme, context *cfg.Config) *HtmlPage {
	log := log.GetLogger()
	log.Trace("NewHtmlPage()")
	if nil == scheme {
		log.Error("nil Scheme, impossible to render HtmlPage!");
		return nil
	}
	return &HtmlPage{
		document:	nil,
		scheme:		scheme,
		context:	context,
	}
}

// Get the rendered document for a given page as specified by the Scheme and context data
func (page *HtmlPage) GetRenderedDocument() string {
	if nil == page.document {
		page.document = page.renderDocument()
	}
	return *page.document
}

// Flatten the Scheme sources, and hydrate with context properties into a final, rendered Document
func (page *HtmlPage) renderDocument() *string {

	// Start our document with the layout content for the scheme
	var tmp *string
	tmp = &page.scheme.GetLayout().GetContent()

	// Dereference layout fragments
	tmp = page.dereferenceFragments(tmp)

	// Inject Scheme stylesheet into the document
	if nil == page.context {
		page.context = cfg.NewConfig()
	}
	page.context.Add("stylesheet", page.scheme.GetStylesheet())

	// TODO: Any string translations needed?

	// Dereference Context/Config
	page.Config.DereferenceLoop(page.Config, DEREFERENCE_MAX_LOOPS)
	page.renderedDocument = page.Config.DereferenceString(page.Document)

	// Return the final, rendered document
	return tmp
}

// Dereference all the Scheme's page Fragments, then Dereference the supplied document against them
func (page *htmlPage) dereferenceFragments(document *string) *string {
	fragments := cfg.NewConfig()
	fragmap := page.scheme.GetFragMap()
	for fragname, fragment := range fragmap {
		// Fragment magic tags are as '%frag:fragment_name%'
		fragments[fmt.Sprintf("frag:%s", fragname)] = fragment.Content
	}
	fullyResolved := fragments.DereferenceLoop(DEREFERENCE_MAX_LOOPS, fragments)
	if ! fullyResolved {
		log.GetLogger().Warn(fmt.Sprintf(
			"Possible incomplete dereferencing of page fragments with %d loops",
			DEREFERENCE_MAX_LOOPS,
		));
	}
	return &fragments.DereferenceString(*document
}

