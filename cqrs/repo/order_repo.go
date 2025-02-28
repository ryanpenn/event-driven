package repo

import (
	"context"
	"errors"
	"event-driven-system/cqrs/model"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
)

// OrderRepository is an interface for order repository
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	DeleteOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context) ([]*model.Order, error)
}

// OrderStatusRepository is an interface for order status repository
type OrderStatusRepository interface {
	GetOrderStatus(ctx context.Context, id string) (string, error)
	UpdateOrderStatus(ctx context.Context, id string, status string) error
}
