# Implementing an Event Driven System in Go

> V2 is an event system based on pub/sub pattern.

## Pub/Sub Pattern

```go
// Subscriber is a function that takes a topic of type T and handles it.
type Subscriber[T any] func(topic T)

// Unsubscriber is an interface that unsubscribes a subscriber.
type Unsubscriber interface {
	// Unsubscribe unsubscribes the subscriber.
	Unsubscribe()
	// Bind binds the Unsubscriber to an AutoCloser Object.
	Bind(AutoCloser)
}

// AutoCloser is an interface that provides a Close method to unsubscribe all subscribers.
type AutoCloser interface {
	// Close closes the AutoCloser and unsubscribes all subscribers.
	Close()
	// Add adds an unsubscriber to the list.
	Add(unsubscriber Unsubscriber)
}
```
Implementation of the pub/sub pattern in [pubsub.go](./pubsub/pubsub.go).

## Implement Events

```go
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
```

Implementation of the event in [event.go](./event/event.go).

## EventBus for events

```go
type Bus interface {
	OnEvent(event EventType, fn pubsub.Subscriber[Event]) pubsub.Unsubscriber
	Trigger(event Event)
}
```
Implementation of the event bus in [bus.go](./event/bus.go).

## Examples

- [Example](./example/main.go)

```bash
go run ./example
```