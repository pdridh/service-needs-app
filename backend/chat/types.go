package chat

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Sender    string             `json:"firstName" bson:"first_name"`
	Receiver  string             `json:"lastName" bson:"last_name"`
	Message   string             `json:"message" bson:"message"` // TODO replace this with something more dynamic (attachments, etc..,etc..)
	Timestamp primitive.DateTime `json:"timestamp" bson:"timestamp"`
}

func NewChatMessage(sender string, receiver string, message string, timestamp time.Time) *ChatMessage {
	return &ChatMessage{
		Sender:    sender,
		Receiver:  receiver,
		Message:   message,
		Timestamp: primitive.NewDateTimeFromTime(timestamp),
	}
}
