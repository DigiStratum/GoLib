package json

type JsonSerializableIfc interface {
	ToJson() (*string, error)
}

type JsonDeserializableIfc interface {
	FromJson(jsonString *string) error
}
