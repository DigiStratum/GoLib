package webui

import(
	log "github.com/DigiStratum/GoLib/Logger"
)

type Scheme struct {
	layout		*Layout
	fragmap		map[string]*Fragment	// key is Fragment.name
	stylesheet	[]*Style
}

// Make a new one of these
func NewScheme() *Scheme {
	log.GetLogger().Trace("NewScheme()")
	return &Scheme{
		layout: nil,
		fragmap: make(map[string]*Fragment),
		stylesheet: make([]*Style, 0),
	}
}

// Set the layout for this scheme
func (scheme *Scheme) SetLayout(layout *Layout) {
	scheme.layout = layout
}

// Get the layout for this scheme
func (scheme *Scheme) GetLayout() *Layout {
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

// Flatten the styles into a single, ordered, stylessheet (CSS)
func (scheme *Scheme) GetStylesheet() string {
	styles := ""
	for style := range scheme.stylesheet {
		styles += "\n" + scheme.stylesheet[style].ToString();
	}
	return styles
}

