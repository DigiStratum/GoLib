package json

/*
Our own take on JsonSerializable.

Note that this is a weak substitute/attempt to "normalize" JSON support for our own library when we
have these available from Go built-in packages:

MarshalJSON() ([]byte, error)		// Equivalent of JsonSerializableIfc
UnmarshalJSON(value []byte) error	// Equivalent of JsonDeserializableIfc

The only difference, really, is []byte vs *string. If we want to seriously support reorienting around
strings, then we probably need to support some sort of reflection/deep recursion mechanism that will
traverse the structures received properly as with Un|MarshallJSON
*/
type JsonSerializableIfc interface {
	ToJson() (*string, error)
}

type JsonDeserializableIfc interface {
	FromJson(jsonString *string) error
}
