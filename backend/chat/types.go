package chat

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Sender    string             `json:"sender" bson:"sender"`
	Receiver  string             `json:"receiver" bson:"receiver"`
	Message   string             `json:"message" bson:"message"` // TODO replace this with something more dynamic (attachments, etc..,etc..)
	Status    MessageStatus      `json:"status" bson:"status"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"created_at"`
}

func NewChatMessage(sender string, receiver string, message string, timestamp time.Time) *ChatMessage {
	return &ChatMessage{
		Sender:    sender,
		Receiver:  receiver,
		Message:   message,
		Timestamp: primitive.NewDateTimeFromTime(timestamp),
	}
}
