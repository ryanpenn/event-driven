package pubsub

import (
	"slices"
	"sync"
)

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

// UnsubscriberFunc is a function that implements the AutoCloser interface.
type UnsubscriberFunc func()

func (f UnsubscriberFunc) Unsubscribe() {
	if f != nil {
		f()
	}
}

func (f UnsubscriberFunc) Bind(ac AutoCloser) {
	if ac != nil {
		ac.Add(f)
	}
}

// AutoUnsub is an object that implements the AutoCloser interface.
type AutoUnsub struct {
	lock          *sync.Mutex
	unsubscribers []Unsubscriber
}

func (au *AutoUnsub) Close() {
	au.lock.Lock()
	defer au.lock.Unlock()

	for _, u := range au.unsubscribers {
		if u != nil {
			u.Unsubscribe()
		}
	}
	au.unsubscribers = au.unsubscribers[:0]
}

func (au *AutoUnsub) Add(unsubscriber Unsubscriber) {
	au.lock.Lock()
	defer au.lock.Unlock()

	au.unsubscribers = append(au.unsubscribers, unsubscriber)
}

var _ AutoCloser = (*AutoUnsub)(nil)

func NewAutoUnsub() *AutoUnsub {
	return &AutoUnsub{
		lock:          &sync.Mutex{},
		unsubscribers: make([]Unsubscriber, 0),
	}
}

// Broker is a pub/sub implementation.
type Broker[T any] struct {
	lock        *sync.RWMutex
	subscribers []*Subscriber[T]
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		lock:        &sync.RWMutex{},
		subscribers: make([]*Subscriber[T], 0),
	}
}

func (b *Broker[T]) Subscribe(subscriber Subscriber[T]) Unsubscriber {
	b.lock.Lock()
	b.subscribers = append(b.subscribers, &subscriber)
	b.lock.Unlock()

	return UnsubscriberFunc(func() {
		if b != nil {
			b.lock.Lock()
			defer b.lock.Unlock()

			b.subscribers = slices.DeleteFunc(b.subscribers, func(s *Subscriber[T]) bool {
				return s == &subscriber
			})
		}
	})
}

func (b *Broker[T]) UnsubscribeAll() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.subscribers = b.subscribers[:0]
}

func (b *Broker[T]) Publish(topic T) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for _, handler := range b.subscribers {
		(*handler)(topic)
	}
}
