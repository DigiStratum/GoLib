package webui

// An HTML Page is a flattened rendering of a Document with all magic tags resolved from the
// supplied Config. The Scheme shall be concatenated into a single block and injected into the
// top of the rendered Document.

// TODO: Deliver schema as a separate CSS file, separately requested (but referenced by Document).
// TODO: Use an external cache (memcache/elasticache?) to make rendered Documents available across
//	 instances using Configuration fingerprinting for differentiation

import(
	lib "github.com/DigiStratum/GoLib"
)

type HtmlPage struct {
	renderedDocument	string
	Document		string
	Config			*lib.Config	// Contextual/Configuration data available for injection into Document
	Scheme			*Scheme
}

// Make a new one of these
func NewHtmlPage() *HtmlPage {
	lib.GetLogger().Trace("NewHtmlPage()")
	return &HtmlPage{
		Document:		"",
		renderedDocument:	nil,
		Scheme:			nil,
		Config:			lib.NewConfig(),
	}
}

func (page *HtmlPage) GetRenderedDocument() string {
	if nil == page.renderedDocument {
		page.renderDocument()
	}
	return page.renderedDocument
}

// Recursively hydrate the source Document into a final, rendered Document
func (page *HtmlPage) renderDocument() {
	// Load Fragments into Config

	// Dereference Config against itself (up to 5 iterations for nested dereferencing)
	// TODO: does this even make sense? Can any validreferences even survive a single pass?
	page.Config.DereferenceAll(page.Config, 5)

	// Dereference the renderedDocument against our Config
	page.renderedDocument = page.Config.DereferenceString(page.Document)
}

