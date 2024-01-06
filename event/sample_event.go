package event

// SampleEventPayload is the data for sampleEvent
type SampleEventPayload struct {
	Name string
}

type sampleEvent struct {
	Event[SampleEventPayload]
}

var SampleEvent sampleEvent

type sampleRefEvent struct {
	Event[*SampleEventPayload]
}

var SampleRefEvent sampleRefEvent
