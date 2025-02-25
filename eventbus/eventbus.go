package eventbus

import (
	"reflect"
	"sync"
)

// EventType is the type of event
type EventType int

// Type implements the Event interface
func (t EventType) Type() EventType {
	return t
}

// Event is the interface for an event
type Event interface {
	// Type is identifies the event
	Type() EventType
}

// EventHandler is the type of event handler
type EventHandler func(event Event)

// EventBus is the interface for an event bus
type EventBus interface {
	// Publish publishes an event to the event bus
	Publish(event Event)
	// Subscribe subscribes an event handler to an event type
	Subscribe(etype EventType, handler EventHandler)
	// Unsubscribe unsubscribes an event handler from an event type
	Unsubscribe(etype EventType, handler EventHandler)
	// UnsubscribeAll unsubscribes all event handlers from an event type
	UnsubscribeAll(etype EventType)
}

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
