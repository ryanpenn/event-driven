package main

import (
	"event-driven-system/event"
	"fmt"
	"time"
)

func init() {
	deleteNotifier := userDeletedNotifier{
		adminEmail: "the.boss@example.com",
		slackHook:  "https://hooks.slack.com/services/delete",
	}

	event.UserDeletedEvent.Register(event.EventID(1), deleteNotifier)
}

type userDeletedNotifier struct {
	adminEmail string
	slackHook  string
}

func (u userDeletedNotifier) notifyAdmin(email string, t time.Time) {
	// send a message to the admin that a user was created
	fmt.Println("notifyAdmin", "deleted", email, t.Format(time.DateTime))
}

func (u userDeletedNotifier) sendToSlack(email string, t time.Time) {
	// send to a slack webhook that a user was created
	fmt.Println("sendToSlack", u.slackHook, email, t.Format(time.DateTime))
}

func (u userDeletedNotifier) Handle(_ event.Sender, payload event.UserDeletedPayload) {
	// Do something with our payload
	u.notifyAdmin(payload.Email, payload.Time)
	u.sendToSlack(payload.Email, payload.Time)
}
