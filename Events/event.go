package events

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
	atiJsonBytes, err := json.Marshal(ati)
	if nil != err { return nil, err }
	atiJsonString := string(atiJsonBytes[:])
	return &atiJsonString, nil
}
