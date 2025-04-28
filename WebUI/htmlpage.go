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
	tv := page.scheme.GetLayout().GetContent()
	tmp := &tv

	// Dereference layout fragments
	tmp = page.dereferenceFragments(tmp)

	// Inject Scheme stylesheet into the document
	if nil == page.context {
		page.context = cfg.NewConfig()
	}
	page.context.Set("stylesheet", page.scheme.GetStylesheet())

	// TODO: Any string translations needed?

	// Dereference Context/Config
	page.context.DereferenceLoop(DEREFERENCE_MAX_LOOPS, page.context)
	page.document = page.context.DereferenceString(*tmp)

	// Return the final, rendered document
	return page.document
}

// Dereference all the Scheme's page Fragments, then Dereference the supplied document against them
func (page *HtmlPage) dereferenceFragments(document *string) *string {
	fragments := cfg.NewConfig()
	fragmap := page.scheme.GetFragMap()
	for fragname, fragment := range fragmap {
		// Fragment magic tags are as '%frag:fragment_name%'
		fragments.Set(fmt.Sprintf("frag:%s", fragname), fragment.Content)
	}
	fullyResolved := fragments.DereferenceLoop(DEREFERENCE_MAX_LOOPS, fragments)
	if ! fullyResolved {
		log.GetLogger().Warn("%s", fmt.Sprintf(
			"Possible incomplete dereferencing of page fragments with %d loops",
			DEREFERENCE_MAX_LOOPS,
		));
	}
	rv := fragments.DereferenceString(*document)
	return rv
}

