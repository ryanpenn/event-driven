package event

import "event-driven-system/eventbus"

type OrderDeletedEvent struct {
	ID string
}

const OrderDeletedEventType eventbus.EventType = 101

// Type implements eventbus.Event.
func (e *OrderDeletedEvent) Type() eventbus.EventType {
	return OrderDeletedEventType
}

var _ eventbus.Event = (*OrderDeletedEvent)(nil)
