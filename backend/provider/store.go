package provider

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProviderStore interface {
	GetProviders(filters bson.M, page int, limit int) ([]bson.M, error)
	InsertProvider(u string, p *Provider) error
}

type mongoProviderStore struct {
	coll *mongo.Collection
}

// Creates a new provider mongo store which uses the given collection
func NewMongoStore(coll *mongo.Collection) *mongoProviderStore {
	return &mongoProviderStore{coll: coll}
}

func (s *mongoProviderStore) GetProviders(filters bson.M, page int, limit int) ([]bson.M, error) {

	skip := (page - 1) * limit
	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))

	cur, err := s.coll.Find(context.TODO(), filters, findOptions)
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

// Given a user's id as a hex string and a provider this function inserts the provider into the collection
func (s *mongoProviderStore) InsertProvider(u string, p *Provider) error {
	p.ID = primitive.NewObjectID()

	uidObj, err := primitive.ObjectIDFromHex(u)
	if err != nil {
		return err
	}
	p.UserID = uidObj

	p.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())

	_, err = s.coll.InsertOne(context.TODO(), p)
	return err
}
