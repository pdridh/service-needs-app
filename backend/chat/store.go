package chat

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	CreateChatMessage(ctx context.Context, c *ChatMessage) error
	GetMessageByID(ctx context.Context, id string) (*ChatMessage, error)
	GetMessagesForChat(ctx context.Context, sender string, receiver string) ([]ChatMessage, error)
	HasMessagedBefore(ctx context.Context, sender string, receiver string) (bool, error)
	UpdateMessageStatus(ctx context.Context, id string, status MessageStatus) error
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
	c.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := s.coll.InsertOne(ctx, c)
	return err
}

func (s *mongoStore) HasMessagedBefore(ctx context.Context, sender string, receiver string) (bool, error) {
	filter := bson.M{
		"sender":   sender,
		"receiver": receiver,
	}

	var c ChatMessage
	err := s.coll.FindOne(ctx, filter).Decode(&c)

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, err
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

func (s *mongoStore) GetMessageByID(ctx context.Context, id string) (*ChatMessage, error) {

	idPrim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": idPrim,
	}

	var c ChatMessage
	err = s.coll.FindOne(ctx, filter).Decode(&c)

	if err == mongo.ErrNoDocuments {
		log.Println("couldnt find coument??")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &c, err
}

func (s *mongoStore) UpdateMessageStatus(ctx context.Context, id string, status MessageStatus) error {

	idPrim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err = s.coll.UpdateByID(ctx, idPrim, update)
	if err != nil {
		return err
	}

	return nil
}
