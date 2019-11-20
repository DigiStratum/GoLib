package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Scheme struct {
	layout		*Layout,
	fragments	map[string]Fragment,	// key is Fragment.name
	styles		map[string]Style,	// key is Style.selector
}

// Make a new one of these
func NewScheme() *Scheme {
	lib.GetLogger().Trace("NewScheme()")
	return &Scheme{
		layout: null,
		fragments: make(map[string]Fragment),
		styles: make(map[string]Style),
	}
}

