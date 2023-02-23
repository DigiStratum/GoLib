package config

type ConfigItemIfc interface {
	GetName() string
}

type configItem struct {
	name		string
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewConfigItem(name string) *configItem {
	return &configItem{
		name:		name,
	}
}

// -------------------------------------------------------------------------------------------------
// ConfigItemIfc
// -------------------------------------------------------------------------------------------------

func (r *configItem) GetName() string {
	return r.name
}

