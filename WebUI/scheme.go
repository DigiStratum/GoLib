package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Scheme struct {
	layout		*Layout,
	fragmap		map[string]*Fragment,	// key is Fragment.name
	stylesheet	[]*Style,
}

// Make a new one of these
func NewScheme() *Scheme {
	lib.GetLogger().Trace("NewScheme()")
	return &Scheme{
		layout: nil,
		fragmap: make(map[string]*Fragment),
		styles: make([]*Style),
	}
}

// Set the layout for this scheme
func (scheme *Scheme) SetLayout(layout *Layout) {
	scheme.layout = layout
}

// Get the layout for this scheme
func (scheme *Scheme) GetLayout() string {
	return scheme.layout
}

// Add a fragment to this scheme (should be referenced either by the layout or another fragment)
func (scheme *Scheme) AddFragment(fragment *Fragment) {
	scheme.fragmap[fragment.Name] = fragment
}

// Get the map of fragments added
func (scheme *Scheme) GetFragMap() map[string]*Fragment {
	return scheme.fragmap
}

// Add a style to this scheme's stylesheet
func (scheme *Scheme) AddStyle(style *Style) {
	scheme.stylesheet = append(scheme.stylesheet, style)
}

