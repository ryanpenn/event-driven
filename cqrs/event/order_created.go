package event

import "event-driven-system/eventbus"

type OrderCreatedEvent struct {
	OrderID string
	Status  string
}

const OrderCreatedEventType eventbus.EventType = 100

// Type implements eventbus.Event.
func (e *OrderCreatedEvent) Type() eventbus.EventType {
	return OrderCreatedEventType
}

var _ eventbus.Event = (*OrderCreatedEvent)(nil)
