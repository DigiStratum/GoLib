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
	Scheme			*Scheme
	Config			*Config
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
	// Initialize from the source Document
	page.renderedDocument = page.Document

	// Dereference Config against itself (up to 5 iterations for nested dereferencing)
	page.Config.DereferenceAll(page.Config, 5)
}

