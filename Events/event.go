package events

import (
	"encoding/json"
)

type EventIfc interface {
}

type Event struct {
	properties		map[string]string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new, generic Event
func NewEvent(properties map[string]string) *Event {
	return &Event{
		properties:	properties,
	}
}

// -------------------------------------------------------------------------------------------------
// GoLib/Data/json/JsonSerializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r Event) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(r.properties)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}
