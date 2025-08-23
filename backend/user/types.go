package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserType string

const (
	UserTypeBusiness UserType = "business"
	UserTypeConsumer UserType = "consumer"
)

// TODO add more fields for services and shit, for now only auth
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	Type      UserType           `json:"type" bson:"type"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}

type UserProfile struct {
	ID      string      `json:"id"`
	Email   string      `json:"email"`
	Type    UserType    `json:"type"`
	Details interface{} `json:"details"`
}
