package query

import (
	"context"
	"event-driven-system/cqrs"
	"event-driven-system/eventbus"
)

type EventQueryService struct {
}

func (EventQueryService) ListEvents(ctx context.Context, streamID string) ([]eventbus.Event, error) {
	return cqrs.EventStore.ReadStream(streamID)
}
