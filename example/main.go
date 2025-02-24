package main

import (
	"event-driven-system/event"
	"event-driven-system/pubsub"
	"fmt"
	"sync"
)

const (
	TestEvent1 event.EventType = iota + 1
	TestEvent2
	ValueChangedEvent
)

type ValueChanged struct {
	Old int
	New int
}

func main() {
	bus := event.NewBus()
	unsub := bus.OnEvent(0, func(e event.Event) {
		fmt.Println(e.Type(), "triggered")
	})
	bus.Trigger(event.EventType(0))
	unsub.Unsubscribe()             // This will unsubscribe the callback function
	bus.Trigger(event.EventType(0)) // This will not trigger the callback function

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		aus := pubsub.NewAutoUnsub() // This will create a new AutoCloser instance
		defer aus.Close()            // This will unsubscribe all the subscribers

		bus.OnEvent(TestEvent1, func(e event.Event) {
			fmt.Println(e.Type(), "triggered")
		}).Bind(aus)
		bus.Trigger(TestEvent1)

		bus.OnEvent(TestEvent2, func(e event.Event) {
			fmt.Println(e.Type(), "triggered")
		}).Bind(aus)
		bus.Trigger(TestEvent2)
	}()

	unsub = bus.OnEvent(ValueChangedEvent, func(e event.Event) {
		if data, ok := e.(event.EventData[ValueChanged]); ok {
			fmt.Println(e.Type(), "triggered ->", data.DataAs().Old, data.DataAs().New)
		}
	})
	bus.Trigger(event.NewEvent(ValueChangedEvent, ValueChanged{1, 2}))
	unsub.Unsubscribe()

	wg.Wait()

	bus.Trigger(TestEvent1) // This will not trigger the callback function
	bus.Trigger(TestEvent2) // This will not trigger the callback function
}
