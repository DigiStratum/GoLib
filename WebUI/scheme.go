package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Scheme struct {
	layout		*Layout,
	fragments	map[string]*Fragment,	// key is Fragment.name
	stylesheet	[]*Style,
}

// Make a new one of these
func NewScheme() *Scheme {
	lib.GetLogger().Trace("NewScheme()")
	return &Scheme{
		layout: nil,
		fragments: make(map[string]*Fragment),
		styles: make([]*Style),
	}
}

// Set the layout for this scheme
func (scheme *Scheme) SetLayout(layout *Layout) {
	scheme.layout = layout
}

// Add a fragment to this scheme (should be referenced either by the layout or another fragment)
func (scheme *Scheme) AddFragment(fragment *Fragment) {
	scheme.fragments[fragment.Name] = fragment
}

// Add a style to this scheme's stylesheet
func (scheme *Scheme) AddStyle(style *Style) {
	scheme.stylesheet = append(scheme.stylesheet, style)


