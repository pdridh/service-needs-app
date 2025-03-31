package provider

import "go.mongodb.org/mongo-driver/bson/primitive"

type Provider struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"-" bson:"user_id"`
	Name        string             `json:"name" bson:"name"`
	Category    string             `json:"category" bson:"category"`
	Location    string             `json:"location" bson:"location"` // TODO change this to somehting better like coords or something idk
	Description string             `json:"description" bson:"description"`
	Verified    bool               `json:"verified" bson:"verified"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	// TODO add other information like available time, documents, profile stuff etc,etc...
}
