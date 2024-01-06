package event

type HandleFunc[T any] func(T)

func (h HandleFunc[T]) Handle(payload T) {
	h(payload)
}

type Event[T any] struct {
	Async    bool
	handlers []interface{ Handle(T) }
}

// RegisterFunc adds an event handler func for this event
func (e *Event[T]) RegisterFunc(fn HandleFunc[T]) {
	e.handlers = append(e.handlers, fn)
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
