package events

type EventProducerIfc interface {
	ProduceEvent(event EventIfc)
}