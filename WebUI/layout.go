package webui

import(
	lib "github.com/DigiStratum/GoLib"
)

type Layout struct {
}

// Make a new one of these
func NewLayout() *Layout {
	lib.GetLogger().Trace("NewLayout()")
	return &Layout{
	}
}

