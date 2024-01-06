# Implementing an Event Driven System in Go

> refer: [event-driven-system](https://stephenafamo.com/blog/posts/implementing-an-event-driven-system-in-go)

## Defining Base Events

```go
// Event [T] is the base event
type Event[T any] struct {
	Async    bool
	handlers []interface{ Handle(T) }
}

// Register adds an event handler for this event
func (e *Event[T]) Register(handler interface{ Handle(T) }) {
	e.handlers = append(e.handlers, handler)
}

// Trigger sends out an event with the payload
func (e *Event[T]) Trigger(payload T) {
	for _, handler := range e.handlers {
		if e.Async {
			go handler.Handle(payload)
		} else {
			handler.Handle(payload)
		}
	}
}
```
- Functions are first-class in Go

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
// SampleEventPayload is the data for sampleEvent
type SampleEventPayload struct {
	Name string
}

// In sampleEvent SampleEventPayload is pass by value
type sampleEvent struct {
	Event[SampleEventPayload]
}

// In sampleRefEvent SampleEventPayload is pass by reference
type sampleRefEvent struct {
	Event[*SampleEventPayload]
}

// Define event instances
var (
    SampleEvent sampleEvent
    SampleRefEvent sampleRefEvent
)
```

## Listening for Events

```go
type sampleNotifier struct {
}

// implement the Handle interface
func (sampleNotifier) Handle(payload event.SampleEventPayload) {
	fmt.Println("sampleNotifier Handle", payload.Name, fmt.Sprintf("%p", &payload))
}

// Register events in init()
func init() {
	notifier := sampleNotifier{}
	event.SampleEvent.Register(notifier)

	event.SampleEvent.RegisterFunc(func(payload event.SampleEventPayload) {
		fmt.Println("FuncHandle", payload.Name, fmt.Sprintf("%p", payload))
	})
}
```

## Triggering Events

```go
func Publish() {
	event.SampleEvent.Trigger(event.SampleEventPayload{
		Name: "Sample",
	})

    // trigger event in goroutine
	event.SampleEvent.Async = true
	event.SampleEvent.Trigger(event.SampleEventPayload{
		Name: "Async Sample",
	})
}
```

## Conclusion

```go
func main() {
	Publish()
	time.Sleep(time.Second) // wait for handling the async event
}
```