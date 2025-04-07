package ws

import (
	"encoding/json"
	"time"
)

type EventCode string

const (
	EventChat  EventCode = "chat"
	EventHello EventCode = "hello"
)

// Actual event that is sent and received by the client.
type Event struct {
	Code    EventCode `json:"code"`
	Payload any       `json:"payload"`
}

// This method marshals the payload and then loads that into v. v is expected to be a ptr, otherwise the changes wont reflect for the caller.
// Returns error if anything failed, nil otherwise.
func (e *Event) ParsePayloadInto(v any) error {
	data, err := json.Marshal(e.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Contains information for the event with the Event itself for the hub.
type EventContext struct {
	Event     Event     `json:"event"`
	Client    *Client   `json:"client"`
	Timestamp time.Time `json:"timestamp"`
}

type EventHelloPayload struct {
	Message string `json:"message"`
}

type EventChatPayload struct {
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
