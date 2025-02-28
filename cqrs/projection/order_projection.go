package projection

import (
	"context"
	"event-driven-system/cqrs"
	"event-driven-system/cqrs/event"
	"event-driven-system/eventbus"
)

type OrderProjectionBuilder struct {
}

// Update implements ProjectionBuilder.
func (*OrderProjectionBuilder) Update(e eventbus.Event) {
	var err error
	switch evt := e.(type) {
	case *event.OrderCreatedEvent:
		err = cqrs.OrderStatusRepository.UpdateOrderStatus(context.Background(), evt.OrderID, "created")
	case *event.OrderDeletedEvent:
		err = cqrs.OrderStatusRepository.UpdateOrderStatus(context.Background(), evt.ID, "deleted")
	}
	if err != nil {
		// handle error
	}
}

func (b *OrderProjectionBuilder) Init() {
	cqrs.EventBus.Subscribe(event.OrderCreatedEventType, b.Update)
	cqrs.EventBus.Subscribe(event.OrderDeletedEventType, b.Update)
}

func (b *OrderProjectionBuilder) Close() {
	cqrs.EventBus.Unsubscribe(event.OrderCreatedEventType, b.Update)
	cqrs.EventBus.Unsubscribe(event.OrderDeletedEventType, b.Update)
}

var _ ProjectionBuilder = (*OrderProjectionBuilder)(nil)

// NewOrderProjectionBuilder creates a new OrderProjectionBuilder.
func NewOrderProjectionBuilder() *OrderProjectionBuilder {
	return &OrderProjectionBuilder{}
}
