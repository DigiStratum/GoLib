package json

type JsonSerializableIfc interface {
	ToJson() (*string, error)
}
