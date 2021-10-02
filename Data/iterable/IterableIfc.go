package iterable

type IterableIfc interface {
	GetIterator() func () interface{}
}
