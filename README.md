# Implementing an Event Driven System in Go

> refer: [event-driven-system](https://stephenafamo.com/blog/posts/implementing-an-event-driven-system-in-go)

## Define Base Events

```go
// Event [T] is the base event
type Event[T any] struct {
	handlers []interface{ Handle(T) }
}

// Register adds an event handler for this event
func (e *Event[T]) Register(handler interface{ Handle(T) }) {
	e.handlers = append(e.handlers, handler)
}

// Trigger sends out an event with the payload
func (e *Event[T]) Trigger(payload T) {
	for _, handler := range e.handlers {
		handler.Handle(payload)
	}
}
```
- Functions are first-class values

```go
// EventHandleFunc
type HandleFunc[T any] func(T)

// implement Handle interface
func (h HandleFunc[T]) Handle(payload T) {
	h(payload)
}

// RegisterFunc adds an event handler func for this event
func (e *Event[T]) RegisterFunc(fn HandleFunc[T]) {
	e.handlers = append(e.handlers, fn)
}
```

## Implement Events

```go
// SamplePayload is the data for sampleEvent
type SamplePayload struct {
	Name string
}

// In sampleEvent SampleEventPayload is pass by value
type sampleEvent struct {
	Event[SamplePayload]
}

// In sampleRefEvent SampleEventPayload is pass by reference
type sampleRefEvent struct {
	Event[*SamplePayload]
}

// Define event instances
var (
    SampleEvent sampleEvent
    SampleRefEvent sampleRefEvent
)
```

## Register handler for events

```go
type sampleNotifier struct {
}

// implement the Handle interface
func (sampleNotifier) Handle(payload event.SamplePayload) {
	fmt.Println("sampleNotifier Handle", payload.Name, fmt.Sprintf("%p", &payload))
}

// Register events in init()
func init() {
	notifier := sampleNotifier{}
	event.SampleEvent.Register(notifier)

	event.SampleEvent.RegisterFunc(func(payload event.SamplePayload) {
		fmt.Println("FuncHandle", payload.Name, fmt.Sprintf("%p", payload))
	})
}
```

## Triggering events

```go
func Publish() {
	event.SampleEvent.Trigger(event.SamplePayload{
		Name: "Sample",
	})
}

func main() {
	Publish()
}
```

## Examples

- [Example](./event/event_test.go)
