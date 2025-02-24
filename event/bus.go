package event

import (
	"event-driven-system/pubsub"
	"sync"
)

type Bus interface {
	OnEvent(event EventType, fn pubsub.Subscriber[Event]) pubsub.Unsubscriber
	Trigger(event Event)
}

type eventBus struct {
	lock *sync.RWMutex
	subs map[EventType]*pubsub.Broker[Event]
}

func (b *eventBus) OnEvent(event EventType, fn pubsub.Subscriber[Event]) pubsub.Unsubscriber {
	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.subs[event]; !ok {
		b.subs[event] = pubsub.NewBroker[Event]()
	}

	return b.subs[event.Type()].Subscribe(fn)
}

func (b *eventBus) Trigger(event Event) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if _, ok := b.subs[event.Type()]; ok {
		b.subs[event.Type()].Publish(event)
	}
}

func NewBus() Bus {
	return &eventBus{
		lock: &sync.RWMutex{},
		subs: make(map[EventType]*pubsub.Broker[Event]),
	}
}
