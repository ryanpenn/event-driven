package eventstore

import "event-driven-system/eventbus"

type EventStore interface {
	AppendToStream(streamID string, events []eventbus.Event) error
	ReadStream(streamID string) ([]eventbus.Event, error)
}
