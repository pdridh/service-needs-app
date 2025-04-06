package ws

import "time"

const (
	EventSendMessage = "EVENT_SEND_MESSAGE"
	EventHello       = "EVENT_HELLO"
)

// Actual event that is sent and received by the client.
type Event struct {
	Code    string `json:"code"`
	Payload any    `json:"payload"`
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
