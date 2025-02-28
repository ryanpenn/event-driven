package eventstore

import (
	"event-driven-system/eventbus"
	"fmt"
)

type InMemoryEventStore struct {
	streams map[string][]eventbus.Event
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{streams: make(map[string][]eventbus.Event)}
}

func (e *InMemoryEventStore) AppendToStream(streamID string, events []eventbus.Event) error {
	fmt.Println("InMemoryEventStore.AppendToStream", streamID)
	e.streams[streamID] = append(e.streams[streamID], events...)
	return nil
}

func (e *InMemoryEventStore) ReadStream(streamID string) ([]eventbus.Event, error) {
	fmt.Println("InMemoryEventStore.ReadStream", len(e.streams[streamID]))
	return e.streams[streamID], nil
}
