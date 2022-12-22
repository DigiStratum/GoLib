package config

type ConfigItemIfc interface {
	GetName() string
	GetValue() interface{}
}

type configItem struct {
	name	string
	value	interface{}
}

// ------------------------------------------------------------------------------------------------
// Factory Functions
// ------------------------------------------------------------------------------------------------

func NewConfigItem(name string, value interface{}) *configItem {
	return &configItem{
		name:		name,
		value:		value,
	}
}


// ------------------------------------------------------------------------------------------------
// ConfigItemIfc Implementation
// ------------------------------------------------------------------------------------------------

func (r *configItem) GetName() string {
	return r.name
}

func (r *configItem) GetValue() interface{} {
	return r.value
}

