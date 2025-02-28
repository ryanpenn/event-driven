package command

import (
	"context"
	"event-driven-system/commandbus"
	"event-driven-system/cqrs"
	"event-driven-system/cqrs/event"
	"event-driven-system/eventbus"
)

type DeleteOrder struct {
	*commandbus.BaseCommand
	orderId string
}

// NewDeleteOrder creates a new DeleteOrder command.
func NewDeleteOrderCommand(orderId string) *DeleteOrder {
	return &DeleteOrder{
		BaseCommand: &commandbus.BaseCommand{},
		orderId:     orderId,
	}
}

// Execute implements commandbus.Command.
// Subtle: this method shadows the method (*BaseCommand).Execute of DeleteOrder.BaseCommand.
func (d *DeleteOrder) Execute(ctx context.Context) error {
	if err := cqrs.OrderRepository.DeleteOrder(ctx, d.orderId); err != nil {
		return err
	}

	evt := &event.OrderDeletedEvent{
		ID: d.orderId,
	}
	cqrs.EventBus.Publish(evt)
	cqrs.EventStore.AppendToStream(d.orderId, []eventbus.Event{evt})

	return d.BaseCommand.Execute(ctx)
}

// Name implements commandbus.Command.
func (d *DeleteOrder) Name() commandbus.CommandName {
	return "DeleteOrder"
}

var _ commandbus.Command = (*DeleteOrder)(nil)
