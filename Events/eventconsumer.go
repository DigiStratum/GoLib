package events

type EventConsumerIfc interface {
	ConsumeEvent(event EventIfc)
}
