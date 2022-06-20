package serializable

// Capture the state of an instance variable/struct as a string that can be stored
type Serializable interface {
	func Serialize() (*string, error)
}

// Restore the state of an instance variable/struct from a string that was stored
type Deserializable interface {
	func Deserialize(data *string) error
}

