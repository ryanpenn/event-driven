package inmem

import (
	"context"
	"event-driven-system/cqrs/model"
	"event-driven-system/cqrs/repo"
	"strconv"
	"time"
)

// OrderRepoInMem is an interface for order repository
type OrderRepoInMem struct {
	orders map[string]*model.Order
}

var _ repo.OrderRepository = (*OrderRepoInMem)(nil)

func NewOrderRepository() *OrderRepoInMem {
	return &OrderRepoInMem{
		orders: make(map[string]*model.Order),
	}
}

func (s *OrderRepoInMem) genOrderID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (s *OrderRepoInMem) CreateOrder(ctx context.Context, order *model.Order) error {
	order.ID = s.genOrderID()
	s.orders[order.ID] = order
	return nil
}

func (s *OrderRepoInMem) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	order, ok := s.orders[id]
	if !ok {
		return nil, repo.ErrOrderNotFound
	}
	return order, nil
}

func (s *OrderRepoInMem) UpdateOrder(ctx context.Context, order *model.Order) error {
	if _, ok := s.orders[order.ID]; !ok {
		return repo.ErrOrderNotFound
	}
	s.orders[order.ID] = order
	return nil
}

func (s *OrderRepoInMem) DeleteOrder(ctx context.Context, id string) error {
	if _, ok := s.orders[id]; !ok {
		return repo.ErrOrderNotFound
	}
	delete(s.orders, id)
	return nil
}

func (s *OrderRepoInMem) ListOrders(ctx context.Context) ([]*model.Order, error) {
	var orders []*model.Order
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders, nil
}
