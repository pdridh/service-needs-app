package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	InsertUser(u *User) error
}

type MongoUserStore struct {
	coll *mongo.Collection
}

// Helper function to get a user by a dynamic field and its value
// Returns a ptr to a User struct or nil if not found
func (s *MongoUserStore) getUserByField(field string, value string) (*User, error) {
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
func (s *MongoUserStore) GetUserByEmail(email string) (*User, error) {
	return s.getUserByField("email", email)
}

// InsertUser inserts the given user into the collection
// InsertUser is also responsible for timestamp
// DOESNT DO ANY CHECKS ASSUMES THAT IT HAS ALRDY PASSED CHECKS
// returns error if any or nil
func (s *MongoUserStore) InsertUser(u *User) error {

	u.ID = primitive.NewObjectID()
	u.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())

	_, err := s.coll.InsertOne(context.TODO(), u)
	return err
}

// Creates a new mongo store with the collection and returns a ptr to it
func NewMongoStore(coll *mongo.Collection) *MongoUserStore {
	return &MongoUserStore{coll: coll}
}
