package review

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	GetReviews(filters bson.M, options *options.FindOptions) ([]Review, error) // Given a business id return all of its reviews
	CreateReview(r *Review) error                                              // Create a document in this collection with the given review struct ptr (also fills other fields for that struct)
}

type mongoStore struct {
	coll *mongo.Collection
}

func NewMongoStore(coll *mongo.Collection) *mongoStore {
	return &mongoStore{
		coll: coll,
	}
}

func (s *mongoStore) GetReviews(filters bson.M, options *options.FindOptions) ([]Review, error) {
	cur, err := s.coll.Find(context.TODO(), filters, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var reviews []Review
	for cur.Next(context.TODO()) {
		var r Review
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *mongoStore) CreateReview(r *Review) error {

	r.ID = primitive.NewObjectID()
	r.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	r.UpdatedAt = r.CreatedAt

	_, err := s.coll.InsertOne(context.TODO(), r)
	return err
}
