package cqrs

import (
	"event-driven-system/commandbus"
	"event-driven-system/cqrs/eventstore"
	"event-driven-system/cqrs/repo/inmem"
	"event-driven-system/eventbus"
)

var (
	// CommandBus is the command bus
	CommandBus = commandbus.NewCommandBus()
	// EventBus is the event bus
	EventBus = eventbus.NewEventBus()
	// EventStore store the events
	EventStore = eventstore.NewInMemoryEventStore()
	// OrderRepository is the order repository
	OrderRepository = inmem.NewOrderRepository()
	//
	OrderStatusRepository = inmem.NewOrderStatusRepository()
)
