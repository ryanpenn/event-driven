package main

import (
	"event-driven-system/eventbus"
	"fmt"
)

type MyEvent struct {
	Name string
	Age  int
}

func (MyEvent) Type() eventbus.EventType {
	return eventbus.EventType(1)
}

func main() {
	bus := eventbus.NewEventBus()
	// subscribe EventType(0) with anonymous handler
	bus.Subscribe(0, func(event eventbus.Event) {
		fmt.Println("anonymous handler received: ", event.Type())
	})
	bus.Publish(eventbus.EventType(0)) // anonymous handler will receive this event
	bus.UnsubscribeAll(0)              // unsubscribe all handlers of EventType(0)
	bus.Publish(eventbus.EventType(0)) // no handler will receive this event

	handler1 := func(event eventbus.Event) {
		myEvent := event.(MyEvent)
		println("handler1 received: ", myEvent.Name, myEvent.Age)
	}
	handler2 := func(event eventbus.Event) {
		myEvent := event.(MyEvent)
		println("handler2 received: ", myEvent.Name, myEvent.Age)
	}
	bus.Subscribe(MyEvent{}.Type(), handler1)   // subscribe event type 1 with handler1
	bus.Subscribe(MyEvent{}.Type(), handler1)   // subscribing twice will work as well
	bus.Subscribe(MyEvent{}.Type(), handler2)   // subscribe event type 1 with handler2
	bus.Publish(MyEvent{Name: "John", Age: 30}) // all handlers will receive this event

	bus.Unsubscribe(MyEvent{}.Type(), handler1)
	bus.Unsubscribe(MyEvent{}.Type(), handler2)
	bus.Publish(MyEvent{Name: "Jane", Age: 25}) // only one handler1 will receive this event

	bus.UnsubscribeAll(MyEvent{}.Type())
	bus.Publish(MyEvent{Name: "Bob", Age: 40}) // no handler will receive this event
}
