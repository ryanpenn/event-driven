package event

// Event is an interface for all events
type Event interface {
	Type() EventType
}

// EventType is a type of event
type EventType int

func (t EventType) Type() EventType {
	return t
}

// EventData is a generic struct for holding event data
type EventData[T any] struct {
	eType EventType
	data  T
}

func (e EventData[T]) Type() EventType {
	return e.eType
}

func (e EventData[T]) Data() any {
	return e.data
}

func (e EventData[T]) DataAs() T {
	return e.data
}

var _ Event = (*EventData[string])(nil)

func NewEvent[T any](eType EventType, data T) EventData[T] {
	return EventData[T]{eType: eType, data: data}
}
