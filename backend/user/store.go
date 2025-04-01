package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(ctx context.Context, u *User) error
}

type mongoStore struct {
	coll *mongo.Collection
}

// Helper function to get a user by a dynamic field and its value
// Returns a ptr to a User struct or nil if not found
func (s *mongoStore) getUserByField(field string, value string) (*User, error) {
	filter := bson.M{field: value}

	var user User
	err := s.coll.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Returns a ptr to a User struct by the given email
// nil if not found and also returns the err
func (s *mongoStore) GetUserByEmail(email string) (*User, error) {
	return s.getUserByField("email", email)
}

// Given a user ptr adds required fields and inserts into the coll using the given context.
// DOESNT DO ANY CHECKS.
// returns error if any or nil.
func (s *mongoStore) CreateUser(ctx context.Context, u *User) error {

	u.ID = primitive.NewObjectID()
	u.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	u.UpdatedAt = u.CreatedAt

	_, err := s.coll.InsertOne(ctx, u)
	return err
}

// Creates a new mongo store with the collection and returns a ptr to it.
func NewMongoStore(coll *mongo.Collection) *mongoStore {
	return &mongoStore{coll: coll}
}
