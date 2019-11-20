package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type HtmlPage struct {
	document		string
	scheme			*Scheme
}

// Make a new one of these
func NewHtmlPage() *HtmlPage {
	lib.GetLogger().Trace("NewHtmlPage()")
	return &HtmlPage{
		document: "",
	}
}

