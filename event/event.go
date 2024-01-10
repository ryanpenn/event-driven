package event

import "fmt"

type (
	// Sender
	Sender interface {
		GetName() string
	}

	// Watcher
	Watcher[T any] interface {
		BeforeTrigger(EventID, Sender, *T)
		AfterTrigger(EventID, Sender, *T)
	}

	// Locker
	Locker interface {
		Lock()
		Unlock()
	}
)

// EventID
type EventID int

// Event function object
type EventFunc[T any] func(sender Sender, payload T)

// implement Event[T any] interface
func (h EventFunc[T]) Handle(sender Sender, payload T) {
	h(sender, payload)
}

// Event [T] base
type Event[T any] struct {
	locker   Locker
	watcher  Watcher[T]
	handlers map[EventID]interface {
		Handle(sender Sender, payload T)
	}
}

func NewEvent[T any]() *Event[T] {
	return &Event[T]{
		handlers: make(map[EventID]interface {
			Handle(sender Sender, payload T)
		}),
	}
}

func (e *Event[T]) SetWatcher(w Watcher[T]) {
	e.watcher = w
}

func (e *Event[T]) SetLocker(l Locker) {
	e.locker = l
}

// RegisterFunc
func (e *Event[T]) RegisterFunc(id EventID, fn EventFunc[T]) error {
	return e.Register(id, fn)
}

// Register event handler
func (e *Event[T]) Register(id EventID, handler interface{ Handle(Sender, T) }) error {
	if handler == nil {
		return fmt.Errorf("handler is nil")
	}

	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	if _, ok := e.handlers[id]; ok {
		return fmt.Errorf("event id %d exists", id)
	}

	e.handlers[id] = handler
	return nil
}

// Unregister
func (e *Event[T]) Unregister(id EventID) {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	delete(e.handlers, id)
}

// Clear all handlers
func (e *Event[T]) Clear() {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	for k := range e.handlers {
		delete(e.handlers, k)
	}
}

// Contains handler
func (e *Event[T]) Contains(id EventID) bool {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	_, ok := e.handlers[id]
	return ok
}

// Count all handlers
func (e *Event[T]) Count() int {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	return len(e.handlers)
}

// Trigger an event with the payload
func (e *Event[T]) Trigger(sender Sender, payload T) {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	for idx, handler := range e.handlers {
		if e.watcher != nil {
			e.watcher.BeforeTrigger(idx, sender, &payload)
		}
		e.doHandle(idx, handler, sender, payload)
	}
}

// TriggerAsync an event with the payload in goroutine
func (e *Event[T]) TriggerAsync(sender Sender, payload T) {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	for idx, handler := range e.handlers {
		if e.watcher != nil {
			e.watcher.BeforeTrigger(idx, sender, &payload)
		}
		go e.doHandle(idx, handler, sender, payload)
	}
}

func (e *Event[T]) doHandle(id EventID, handler interface{ Handle(Sender, T) }, sender Sender, payload T) {
	defer func() {
		if e.watcher != nil {
			e.watcher.AfterTrigger(id, sender, &payload)
		}
	}()

	handler.Handle(sender, payload)
}
