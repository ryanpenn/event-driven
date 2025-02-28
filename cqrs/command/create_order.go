package command

import (
	"context"
	"event-driven-system/commandbus"
	"event-driven-system/cqrs"
	"event-driven-system/cqrs/event"
	"event-driven-system/cqrs/model"
	"event-driven-system/eventbus"
)

// CreateOrderCommand is a command to create an order.
type CreateOrderCommand struct {
	*commandbus.BaseCommand
	order *model.Order
}

var _ commandbus.Command = (*CreateOrderCommand)(nil)

func NewCreateOrderCommand(order *model.Order) *CreateOrderCommand {
	return &CreateOrderCommand{
		BaseCommand: &commandbus.BaseCommand{},
		order:       order,
	}
}

func (c *CreateOrderCommand) Name() commandbus.CommandName {
	return "CreateOrder"
}

func (c *CreateOrderCommand) Execute(ctx context.Context) error {
	if err := cqrs.OrderRepository.CreateOrder(ctx, c.order); err != nil {
		return err
	}

	evt := &event.OrderCreatedEvent{
		OrderID: c.order.ID,
		Status:  c.order.Status,
	}
	cqrs.EventBus.Publish(evt)
	cqrs.EventStore.AppendToStream(c.order.ID, []eventbus.Event{evt})

	return c.BaseCommand.Execute(ctx)
}
