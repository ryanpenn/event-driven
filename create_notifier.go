package main

import (
	"event-driven-system/event"
	"fmt"
	"time"
)

func init() {
	createNotifier := userCreatedNotifier{
		adminEmail: "the.boss@example.com",
		slackHook:  "https://hooks.slack.com/services/create",
	}

	event.UserCreatedEvent.Register(event.EventID(1), createNotifier)
}

type userCreatedNotifier struct {
	adminEmail string
	slackHook  string
}

func (u userCreatedNotifier) notifyAdmin(email string, t time.Time) {
	// send a message to the admin that a user was created
	fmt.Println("notifyAdmin", "created", email, t.Format(time.DateTime))
}

func (u userCreatedNotifier) sendToSlack(email string, t time.Time) {
	// send to a slack webhook that a user was created
	fmt.Println("sendToSlack", u.slackHook, email, t.Format(time.DateTime))
}

func (u userCreatedNotifier) Handle(_ event.Sender, payload *event.UserCreatedPayload) {
	// Do something with our payload
	u.notifyAdmin(payload.Email, payload.Time)
	u.sendToSlack(payload.Email, payload.Time)
}
