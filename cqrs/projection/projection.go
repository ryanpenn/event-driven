package projection

import "event-driven-system/eventbus"

type ProjectionBuilder interface {
	Update(event eventbus.Event)
}
