package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Fragment struct {
	name		string,
	content		string,
}

// Make a new one of these
func NewFragment() *Fragment {
	lib.GetLogger().Trace("NewFragment()")
	return &Fragment{
		name: "",
		content: "",
	}
}

func (fragment *Fragment) SetName(name string) {
	fragment.name = name
}

func (fragment *Fragment) SetContent(content string) {
	fragment.content = content
}

