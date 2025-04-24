package consumer

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	CreateConsumer(ctx context.Context, c *Consumer) error
	GetConsumerByID(ctx context.Context, id string) (*Consumer, error)
}

type mongoStore struct {
	coll *mongo.Collection
}

// Creates a new mongo store with the collection and returns a ptr to it
func NewMongoStore(coll *mongo.Collection) *mongoStore {
	return &mongoStore{coll: coll}
}

// Given a consumer struct ptr adds required fields and inserts into the collection using the given context
// DOESNT DO ANY CHECKS.
func (s *mongoStore) CreateConsumer(ctx context.Context, c *Consumer) error {
	c.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	c.UpdatedAt = c.CreatedAt

	_, err := s.coll.InsertOne(ctx, c)
	return err
}

func (s *mongoStore) GetConsumerByID(ctx context.Context, id string) (*Consumer, error) {
	idPrim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var c Consumer
	filter := bson.M{"_id": idPrim}
	err = s.coll.FindOne(ctx, filter).Decode(&c)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}
