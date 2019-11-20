package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Style struct {
	selector	string,
	definition	string,
}

// Make a new one of these
func NewStyle() *Style {
	lib.GetLogger().Trace("NewStyle()")
	return &Style{
		selector: "",
		definition: "",
	}
}

func (style *Style) SetSelector(selector string) {
	style.selector = selector
}

func (style *Style) SetDefinition(definition string) {
	style.definition = definition
}

