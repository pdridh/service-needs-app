package chat

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	CreateChatMessage(ctx context.Context, c *ChatMessage) error
	GetMessagesForChat(ctx context.Context, sender string, receiver string) ([]ChatMessage, error)
}

type mongoStore struct {
	coll *mongo.Collection
}

// Creates a new mongo store with the collection and returns a ptr to it
func NewMongoStore(coll *mongo.Collection) *mongoStore {
	return &mongoStore{coll: coll}
}

func (s *mongoStore) CreateChatMessage(ctx context.Context, c *ChatMessage) error {
	c.ID = primitive.NewObjectID()

	_, err := s.coll.InsertOne(ctx, c)
	return err
}

func (s *mongoStore) GetMessagesForChat(ctx context.Context, sender string, receiver string) ([]ChatMessage, error) {
	// For now send all messages between two
	filter := bson.M{
		"$or": []bson.M{
			{"sender": sender, "receiver": receiver},
			{"sender": receiver, "receiver": sender},
		},
	}

	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []ChatMessage
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
