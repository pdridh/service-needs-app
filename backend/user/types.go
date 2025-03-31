package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO add more fields for services and shit, for now only auth
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}
