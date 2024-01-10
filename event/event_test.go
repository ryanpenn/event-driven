package event_test

import (
	"event-driven-system/event"
	"fmt"
	"sync"
	"testing"
	"time"
)

// go test -timeout 30s -run ^TestSimpleEvent$ event-driven-system/event -v
func TestSimpleEvent(t *testing.T) {
	// Payload
	type samplePayload struct {
		Value int
	}

	// Event
	var evt = struct {
		*event.Event[samplePayload]
	}{
		Event: event.NewEvent[samplePayload](),
	}

	t.Run("simple1", func(t *testing.T) {
		// register event handler
		evt.RegisterFunc(event.EventID(1), func(sender event.Sender, payload samplePayload) {
			fmt.Println("payload value", payload.Value)
		})

		// event arg
		payload := samplePayload{
			Value: 10,
		}
		// trigger event
		evt.Trigger(nil, payload)
	})

	t.Run("simple2", func(t *testing.T) {
		evt.Trigger(nil, samplePayload{
			Value: 1,
		})
		// unregister event handler
		evt.Unregister(event.EventID(1))

		// trigger event again, but no handler
		evt.Trigger(nil, samplePayload{
			Value: 2,
		})
	})

}

// -----------------------------------------------------
// advanced usage

type (
	// Payload
	samplePayload struct {
		Value int
	}
	// Sender
	sampleSender struct {
	}
	// Event
	sampleAsyncEvent struct {
		sync.Mutex
		wg sync.WaitGroup
		*event.Event[samplePayload]
	}
)

func (*sampleSender) GetName() string {
	return "sample"
}

func (e *sampleAsyncEvent) Wait() {
	e.wg.Wait()
}

func (e *sampleAsyncEvent) BeforeTrigger(id event.EventID, sender event.Sender, payload *samplePayload) {
	e.wg.Add(1)

	fmt.Printf("BeforeTrigger => id:%d sender:%s payload.Value:%d\n", id, sender.GetName(), payload.Value)
}

func (e *sampleAsyncEvent) AfterTrigger(id event.EventID, sender event.Sender, payload *samplePayload) {
	e.wg.Done()

	fmt.Printf("AfterTrigger => id:%d sender:%s payload.Value:%d\n", id, sender.GetName(), payload.Value)
}

var (
	_ event.Sender                 = (*sampleSender)(nil)
	_ event.Watcher[samplePayload] = (*sampleAsyncEvent)(nil)
	_ event.Locker                 = (*sampleAsyncEvent)(nil)
)

func newAsyncEvent() *sampleAsyncEvent {
	evt := &sampleAsyncEvent{
		Event: event.NewEvent[samplePayload](),
	}
	evt.SetWatcher(evt)
	evt.SetLocker(evt)
	return evt
}

// go test -timeout 30s -run ^TestAsyncEvent$ event-driven-system/event -v
func TestAsyncEvent(t *testing.T) {
	t.Run("trigger async", func(t *testing.T) {
		asyncEvent := newAsyncEvent()

		asyncEvent.RegisterFunc(event.EventID(1), func(sender event.Sender, payload samplePayload) {
			time.Sleep(time.Second * 1)
			fmt.Println("[H1] payload value", payload.Value, "from", sender.GetName())
		})
		asyncEvent.RegisterFunc(event.EventID(2), func(sender event.Sender, payload samplePayload) {
			time.Sleep(time.Second * 2)
			fmt.Println("[H2] payload value", payload.Value, "from", sender.GetName())
		})
		asyncEvent.RegisterFunc(event.EventID(3), func(sender event.Sender, payload samplePayload) {
			time.Sleep(time.Second * 3)
			fmt.Println("[H3] payload value", payload.Value, "from", sender.GetName())
		})

		asyncEvent.TriggerAsync(&sampleSender{}, samplePayload{
			Value: 100,
		})
		asyncEvent.Wait()
	})

	t.Run("thread-safe event", func(t *testing.T) {
		asyncEvent := newAsyncEvent()

		for i := 0; i < 30; i++ {
			asyncEvent.RegisterFunc(event.EventID(i), func(sender event.Sender, payload samplePayload) {
				fmt.Println("handle...", i, payload.Value)
				asyncEvent.Unregister(event.EventID(payload.Value))
				fmt.Println("Unregister ->", payload.Value)
			})

			asyncEvent.TriggerAsync(&sampleSender{}, samplePayload{
				Value: i,
			})

			if i%5 == 0 {
				asyncEvent.Clear()
			}
		}
		asyncEvent.Wait()

		fmt.Println("event handles count", asyncEvent.Count())
	})
}
