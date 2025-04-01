package consumer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Consumer struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserID    primitive.ObjectID `json:"-" bson:"user_id"`
	FirstName string             `json:"firstName" bson:"first_name"`
	LastName  string             `json:"lastName" bson:"last_name"`
	Verified  bool               `json:"verified" bson:"verified"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}
