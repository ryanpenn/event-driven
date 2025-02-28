package inmem

import (
	"context"
	"event-driven-system/cqrs/repo"
)

type OrderStatusRepoInMem struct {
	orderStatus map[string]string
}

// GetOrderStatus implements repo.OrderStatusRepository.
func (o *OrderStatusRepoInMem) GetOrderStatus(ctx context.Context, id string) (string, error) {
	if status, ok := o.orderStatus[id]; ok {
		return status, nil
	}
	return "", repo.ErrOrderNotFound
}

// UpdateOrderStatus implements repo.OrderStatusRepository.
func (o *OrderStatusRepoInMem) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	o.orderStatus[id] = status
	return nil
}

var _ repo.OrderStatusRepository = (*OrderStatusRepoInMem)(nil)

func NewOrderStatusRepository() *OrderStatusRepoInMem {
	return &OrderStatusRepoInMem{
		orderStatus: make(map[string]string),
	}
}
