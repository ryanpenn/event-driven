# Implementing an Event Driven System in Go

> V3 is an eventbus implementation for an event driven system.

## Event Interface

```go
type EventType int

func (t EventType) Type() EventType {
	return t
}

type Event interface {
	Type() EventType
}

type EventHandler func(event Event)

type EventBus interface {
	Publish(event Event)
	Subscribe(etype EventType, handler EventHandler)
	Unsubscribe(etype EventType, handler EventHandler)
	UnsubscribeAll(etype EventType)
}
```

## EventBus Implement

```go
type eventbus struct {
	lock     *sync.RWMutex
	handlers map[EventType]map[uintptr][]EventHandler
}

func NewEventBus() EventBus {
	return &eventbus{
		lock:     &sync.RWMutex{},
		handlers: map[EventType]map[uintptr][]EventHandler{},
	}
}

func (eb *eventbus) Publish(event Event) {
	eb.lock.RLock()
	defer eb.lock.RUnlock()

	handlers := eb.handlers[event.Type()]
	for _, handler := range handlers {
		for _, fn := range handler {
			fn(event)
		}
	}
}

func (eb *eventbus) Subscribe(etype EventType, handler EventHandler) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	if _, ok := eb.handlers[etype]; !ok {
		eb.handlers[etype] = map[uintptr][]EventHandler{}
	}

	ptr := reflect.ValueOf(handler).Pointer()
	eb.handlers[etype][ptr] = append(eb.handlers[etype][ptr], handler)
}

func (eb *eventbus) Unsubscribe(etype EventType, handler EventHandler) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	if _, ok := eb.handlers[etype]; !ok {
		return
	}

	ptr := reflect.ValueOf(handler).Pointer()
	if _, ok := eb.handlers[etype][ptr]; !ok {
		return
	}

	l := len(eb.handlers[etype][ptr])
	if l > 0 {
		eb.handlers[etype][ptr] = eb.handlers[etype][ptr][:l-1]
	}

	if len(eb.handlers[etype][ptr]) == 0 {
		delete(eb.handlers[etype], ptr)
	}
}

func (eb *eventbus) UnsubscribeAll(etype EventType) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	if _, ok := eb.handlers[etype]; !ok {
		return
	}

	for ptr := range eb.handlers[etype] {
		delete(eb.handlers[etype], ptr)
	}
}
```

all code in [eventbus.go](./eventbus/eventbus.go).

## Examples

- [Example](./example/main.go)

```bash
go run ./example
```