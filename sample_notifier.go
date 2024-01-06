package main

import (
	"event-driven-system/event"
	"fmt"
)

func init() {
	notifier := sampleNotifier{}
	event.SampleEvent.Register(notifier)
	event.SampleEvent.RegisterFunc(MyFuncHandle)

	event.SampleRefEvent.RegisterFunc(func(payload *event.SampleEventPayload) {
		fmt.Println("SampleRefEvent FuncHandle1", payload.Name, fmt.Sprintf("%p", payload))
	})
	event.SampleRefEvent.RegisterFunc(func(payload *event.SampleEventPayload) {
		fmt.Println("SampleRefEvent FuncHandle2", payload.Name, fmt.Sprintf("%p", payload))
	})
}

type sampleNotifier struct {
}

func (sampleNotifier) Handle(payload event.SampleEventPayload) {
	fmt.Println("sampleNotifier Handle", payload.Name, fmt.Sprintf("%p", &payload))
}

func MyFuncHandle(payload event.SampleEventPayload) {
	fmt.Println("MyFuncHandle", payload.Name, fmt.Sprintf("%p", &payload))
}
