package ws

import (
	"encoding/json"
	"time"
)

type EventCode string

const (
	EventChat     EventCode = "chat"
	EventHello    EventCode = "hello"
	EventChatSeen EventCode = "chatseen"
)

// Actual event that is sent and received by the client.
type Event struct {
	Code    EventCode `json:"code"`
	Payload any       `json:"payload"`
}

// This method marshals the payload and then loads that into v. v is expected to be a ptr, otherwise the changes wont reflect for the caller.
// Returns error if anything failed, nil otherwise. Can also be used to check if the payload for the event is valid.
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

type EventSeenPayload struct {
	MessageID string `json:"messageID"`
}
