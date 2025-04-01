package review

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	BusinessID primitive.ObjectID `json:"businessID" bson:"business_id"`
	ConsumerID primitive.ObjectID `json:"consumerID" bson:"consumer_id"`
	Comment    string             `json:"comment" bson:"comment"`
	Rating     float32            `json:"rating" bson:"rating"`
	CreatedAt  primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt  primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}
