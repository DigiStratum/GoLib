package webui

// An HTML Page is a flattened rendering of a Document with all magic tags resolved from the
// supplied Config. The Scheme shall be concatenated into a single block and injected into the
// top of the rendered Document.

// TODO: Deliver Scheme.stylesheet as a separate CSS file, separately requested (but referenced by Document).
// TODO: Use an external cache (memcache/elasticache?) to make rendered Documents available across
//	 instances using Configuration fingerprinting for differentiation

import(
	lib "github.com/DigiStratum/GoLib"
)

type HtmlPage struct {
	config			*lib.Config	// Contextual/Configuration data available for injection into Document
	scheme			*Scheme
	document		string
}

// Make a new one of these
func NewHtmlPage() *HtmlPage {
	lib.GetLogger().Trace("NewHtmlPage()")
	return &HtmlPage{
		document:		nil,
		scheme:			nil,
		config:			lib.NewConfig(),
	}
}

func (page *HtmlPage) GetRenderedDocument() string {
	if nil == page.document {
		page.renderDocument()
	}
	return page.document
}

// Recursively hydrate the source Document into a final, rendered Document
func (page *HtmlPage) renderDocument() {

	// Start our document with the layout content for the scheme
	tmp := page.scheme.GetLayout().GetContent()

	// Dereference layout fragments
	fragments := lib.NewConfig()
	fragmap := scheme.GetFragMap()
	for fragname, fragment := range fragmap {
		fragments[fragname] = fragment.Content
	}
	fragments.DereferenceLoop(fragments, 5)
	ttmp = fragments.DereferenceString(tmp)
	tmp = *ttmp

	// TODO: Inject stylesheet into the document
	// TODO: Any string translations needed?

	// Dereference Config
	page.Config.DereferenceLoop(page.Config, 5)
	page.renderedDocument = page.Config.DereferenceString(page.Document)

	// Capture the final, rendered document
	page.document = tmp
}

