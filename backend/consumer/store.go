package consumer

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	CreateConsumer(ctx context.Context, c *Consumer) error
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
