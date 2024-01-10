package user

import (
	"event-driven-system/event"
	"time"
)

func CreateUser() {
	// ...
	event.UserCreatedEvent.Trigger(nil, &event.UserCreatedPayload{
		Email: "new.user@example.com",
		Time:  time.Now(),
	})
	// ...
}

func DeleteUser() {
	// ...
	event.UserDeletedEvent.Trigger(nil, event.UserDeletedPayload{
		Email: "deleted.user@example.com",
		Time:  time.Now(),
	})
	// ...
}
