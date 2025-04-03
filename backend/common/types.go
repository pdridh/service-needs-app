package common

type GeoLocation struct {
	Type        string    `json:"-" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}
