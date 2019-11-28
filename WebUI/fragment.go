package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Fragment struct {
	Name		string,
	Content		string,
}

// Make a new one of these
func NewFragment() *Fragment {
	lib.GetLogger().Trace("NewFragment()")
	return &Fragment{
		Name: "",
		Content: "",
	}
}

