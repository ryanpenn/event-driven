package command

import (
	"context"
	"event-driven-system/commandbus"
	"event-driven-system/cqrs"
	"event-driven-system/cqrs/model"
)

type UpdateOrderCommand struct {
	*commandbus.BaseCommand
	order *model.Order
}

// NewUpdateOrderCommand creates a new UpdateOrderCommand.
func NewUpdateOrderCommand(order *model.Order) *UpdateOrderCommand {
	return &UpdateOrderCommand{
		BaseCommand: &commandbus.BaseCommand{},
		order:       order,
	}
}

// Execute implements commandbus.Command.
// Subtle: this method shadows the method (*BaseCommand).Execute of UpdateOrderCommand.BaseCommand.
func (u *UpdateOrderCommand) Execute(ctx context.Context) error {
	// if err := u.order.Validate(); err != nil {
	// 	return err
	// }
	if err := cqrs.OrderRepository.UpdateOrder(ctx, u.order); err != nil {
		return err
	}

	return u.BaseCommand.Execute(ctx)
}

// Name implements commandbus.Command.
func (u *UpdateOrderCommand) Name() commandbus.CommandName {
	return "UpdateOrder"
}

var _ commandbus.Command = (*UpdateOrderCommand)(nil)
