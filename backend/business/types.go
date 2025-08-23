package business

import (
	"github.com/pdridh/service-needs-app/backend/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents a business entity
type Business struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Category    string              `json:"category" bson:"category"`
	Location    common.GeoJSONPoint `json:"location" bson:"location"`
	Description string              `json:"description" bson:"description"`
	Verified    bool                `json:"verified" bson:"verified"`
	CreatedAt   primitive.DateTime  `json:"createdAt" bson:"created_at"`
	UpdatedAt   primitive.DateTime  `json:"updatedAt" bson:"updated_at"`
	Distance    *float64            `json:"distance,omitempty" bson:"distance,omitempty"`
	AvgRating   *float64            `json:"averageRating,omitempty" bson:"average_rating,omitempty"`
	ReviewCount int                 `json:"reviewCount" bson:"review_count"`
	// TODO add other information like available time, documents, profile stuff etc,etc...
}

type BusinessDetails struct {
	Name        string              `json:"name" bson:"name"`
	Category    string              `json:"category" bson:"category"`
	Location    common.GeoJSONPoint `json:"location" bson:"location"`
	Description string              `json:"description" bson:"description"`
	Verified    bool                `json:"verified" bson:"verified"`
	Distance    *float64            `json:"distance,omitempty" bson:"distance,omitempty"`
	AvgRating   *float64            `json:"averageRating,omitempty" bson:"average_rating,omitempty"`
	ReviewCount int                 `json:"reviewCount" bson:"review_count"`
}

// QueryOptions represents all possible query parameters for business listing
type QueryOptions struct {
	Page      int64          `json:"page"`
	PageSize  int64          `json:"pageSize"`
	SortBy    string         `json:"sortBy"`
	SortOrder string         `json:"sortOrder"` // asc or desc
	Search    string         `json:"search"`
	Filters   map[string]any `json:"filters"`
	Longitude float64        `json:"longitude"`
	Latitude  float64        `json:"latitude"`
	MaxDist   float64        `json:"maxDist"` // in meters
}

type PaginationMetadata struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
}
