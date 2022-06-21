package serializable

/*
General support for de|serialization
*/

// Capture the state of an instance variable/struct as a string that can be stored
type SerializableIfc interface {
	Serialize() (*string, error)
}

// Restore the state of an instance variable/struct from a string that was stored
type DeserializableIfc interface {
	Deserialize(data *string) error
}

