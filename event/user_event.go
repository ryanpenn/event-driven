package event

import "time"

type (
	// UserCreatedPayload is the data for when a user is created
	UserCreatedPayload struct {
		Email string
		Time  time.Time
	}

	// UserDeletedPayload
	UserDeletedPayload struct {
		Email string
		Time  time.Time
	}
)

var (
	// UserCreatedEvent
	UserCreatedEvent = struct {
		*Event[*UserCreatedPayload] // payload pass by pointer
	}{
		Event: NewEvent[*UserCreatedPayload](),
	}

	// UserDeletedEvent
	UserDeletedEvent = struct {
		*Event[UserDeletedPayload] // payload pass by value
	}{
		Event: NewEvent[UserDeletedPayload](),
	}
)
