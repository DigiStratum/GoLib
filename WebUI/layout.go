package webui

import(
	log "github.com/DigiStratum/GoLib/Logger"
)

type Layout struct {
	content			string
}

// Make a new one of these
func NewLayout() *Layout {
	log.GetLogger().Trace("NewLayout()")
	return &Layout{
		content: "",
	}
}

// Set the layout content
func (layout *Layout) SetContent(content string) {
	layout.content = content
}

// Get the layout content
func (layout *Layout) GetContent() string {
	return layout.content
}

