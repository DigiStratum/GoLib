package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Style struct {
	selector	string
	definition	string
	flattened	*string
}

// Make a new one of these
func NewStyle(selector, definition string) *Style {
	lib.GetLogger().Trace("NewStyle()")
	return &Style{
		selector: selector,
		definition: definition,
	}
}

// Flatten into a CSS string ready to use in an HtmlPage
func (style *Style) ToString() string {
	if nil == style.flattened {
		style.flattened = &fmt.Sprintf("%s {\n%s\n}\n", style.selector, style.definition)
	}
	return *style.flattened
}

