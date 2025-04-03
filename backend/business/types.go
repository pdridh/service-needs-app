package business

import (
	"github.com/pdridh/service-needs-app/backend/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents a business entity
type Business struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Category    string             `json:"category" bson:"category"`
	Location    common.GeoLocation `json:"location" bson:"location"` // TODO change this to somehting better like coords or something idk
	Description string             `json:"description" bson:"description"`
	Verified    bool               `json:"verified" bson:"verified"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updated_at"`
	// TODO add other information like available time, documents, profile stuff etc,etc...
}
