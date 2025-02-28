package query

import (
	"context"
	"event-driven-system/cqrs"
	"event-driven-system/cqrs/model"
)

type OrderQueryService struct {
}

func (s OrderQueryService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return cqrs.OrderRepository.GetOrder(ctx, id)
}

func (s OrderQueryService) ListOrders(ctx context.Context) ([]*model.Order, error) {
	return cqrs.OrderRepository.ListOrders(ctx)
}

func (s OrderQueryService) GetOrderStatus(ctx context.Context, id string) (string, error) {
	order, err := cqrs.OrderRepository.GetOrder(ctx, id)
	if err != nil {
		return "", err
	}
	return order.Status, nil
}
