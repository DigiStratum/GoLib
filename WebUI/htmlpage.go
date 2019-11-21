package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type HtmlPage struct {
	document		string
	scheme			*Scheme
	dictionary		map[string]string
	renderedDocument	string
}

// Make a new one of these
func NewHtmlPage() *HtmlPage {
	lib.GetLogger().Trace("NewHtmlPage()")
	return &HtmlPage{
		document: "",
		scheme: nil,
		dictionary: make(map[string]string),
		renderedDocument: nil,
	}
}

func (page *HtmlPage) GetRenderedDocument() string {
	if nil == page.renderedDocument {
		page.renderDocument()
	}
	return page.renderedDocument
}

func (page *HtmlPage) renderDocument() {
	page.renderedDocument = document
}

