package user

import (
	"event-driven-system/event"
	"time"
)

func CreateUser() {
	// ...
	event.UserCreated.Trigger(event.UserCreatedPayload{
		Email: "new.user@example.com",
		Time:  time.Now(),
	})
	// ...
}

func DeleteUser() {
	// ...
	event.UserDeleted.Trigger(event.UserDeletedPayload{
		Email: "deleted.user@example.com",
		Time:  time.Now(),
	})
	// ...
}

func SamplePublish() {
	event.SampleEvent.Trigger(event.SampleEventPayload{
		Name: "Sample",
	})

	event.SampleEvent.Async = true
	event.SampleEvent.Trigger(event.SampleEventPayload{
		Name: "Async Sample",
	})

	event.SampleRefEvent.Trigger(&event.SampleEventPayload{
		Name: "Ref Sample",
	})
}
