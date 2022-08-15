package webui

import(
	log "github.com/DigiStratum/GoLib/Logger"
)

type Fragment struct {
	Name		string,
	Content		string,
}

// Make a new one of these
func NewFragment() *Fragment {
	log.GetLogger().Trace("NewFragment()")
	return &Fragment{
		Name: "",
		Content: "",
	}
}

