package business

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	GetBusinesses(filters bson.M, options *options.FindOptions) ([]bson.M, error)
	GetBusinessByID(id string) (*Business, error)
	CreateBusiness(ctx context.Context, b *Business) error
}

type mongoStore struct {
	coll *mongo.Collection
}

// Creates a new provider mongo store which uses the given collection
func NewMongoStore(coll *mongo.Collection) *mongoStore {
	return &mongoStore{coll: coll}
}

func (s *mongoStore) GetBusinesses(filters bson.M, options *options.FindOptions) ([]bson.M, error) {
	cur, err := s.coll.Find(context.TODO(), filters, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var results []bson.M
	if err := cur.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *mongoStore) GetBusinessByID(id string) (*Business, error) {
	idPrim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var business Business
	filter := bson.M{"_id": idPrim}
	err = s.coll.FindOne(context.TODO(), filter).Decode(&business)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &business, nil
}

// Inserts the given business into the collection using the given ctx.
func (s *mongoStore) CreateBusiness(ctx context.Context, b *Business) error {
	b.ID = primitive.NewObjectID()
	b.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	b.UpdatedAt = b.CreatedAt

	_, err := s.coll.InsertOne(ctx, b)
	return err
}
