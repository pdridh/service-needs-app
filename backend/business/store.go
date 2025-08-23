package business

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	GetBusinesses(ctx context.Context, options QueryOptions) ([]Business, int64, error)
	GetBusinessByID(id string) (*Business, error)
	CreateBusiness(ctx context.Context, b *Business) error
}

type mongoStore struct {
	coll *mongo.Collection
}

// Creates a new provider mongo store which uses the given collection
func NewMongoStore(coll *mongo.Collection) *mongoStore {
	return &mongoStore{coll: coll}
}

func (s *mongoStore) GetBusinesses(ctx context.Context, options QueryOptions) ([]Business, int64, error) {
	// This shit is so horrible i cannot even start, making me regret using mongo ong
	pipeline := mongo.Pipeline{}

	// Add geolocation search if coordinates are provided
	if options.Longitude != 0 || options.Latitude != 0 {
		maxDistance := options.MaxDist
		if maxDistance <= 0 {
			maxDistance = 5000 // Default 5km
		}

		geoNearStage := bson.D{
			{Key: "$geoNear", Value: bson.D{
				{Key: "near", Value: bson.D{
					{Key: "type", Value: "Point"},
					{Key: "coordinates", Value: []float64{options.Longitude, options.Latitude}},
				}},
				{Key: "distanceField", Value: "distance"},
				{Key: "maxDistance", Value: maxDistance},
				{Key: "spherical", Value: true},
			}},
		}
		pipeline = append(pipeline, geoNearStage)
	}

	if options.Search != "" {
		// Create OR conditions for multiple fields
		searchFields := []string{"name", "description", "category"}
		orConditions := bson.A{}

		for _, field := range searchFields {
			orConditions = append(orConditions,
				bson.D{{Key: field, Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: options.Search, Options: "i"}}}}})
		}

		// Add $match stage with $or condition for search
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: orConditions}}}})
	}

	// Add filtering stage if filters are provided
	if len(options.Filters) > 0 {
		filterConditions := bson.D{}
		for field, value := range options.Filters {
			filterConditions = append(filterConditions, bson.E{Key: field, Value: value})
		}
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filterConditions}})
	}

	// Add lookup stage to join with reviews
	pipeline = append(pipeline, bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "reviews"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "business_id"},
			{Key: "as", Value: "reviews"},
		}},
	})

	// Calculate average rating
	pipeline = append(pipeline, bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "average_rating", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$eq", Value: bson.A{
						bson.D{{Key: "$size", Value: "$reviews"}},
						0,
					}}},
					nil,
					bson.D{{Key: "$avg", Value: "$reviews.rating"}},
				}},
			}},
			{Key: "review_count", Value: bson.D{
				{Key: "$size", Value: "$reviews"},
			}},
		}},
	})

	// Get total count for pagination
	countPipeline := make(mongo.Pipeline, len(pipeline))
	copy(countPipeline, pipeline)
	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "total"}})

	// Get count results
	countCursor, err := s.coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}

	// Set default total count
	var totalCount int64 = 0
	if len(countResult) > 0 {
		totalCount = int64(countResult[0]["total"].(int32))
	}

	// Add sorting stage if sort field is provided
	if options.SortBy != "" {
		sortOrder := 1 // Default ascending
		if options.SortOrder == "desc" {
			sortOrder = -1
		}
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: options.SortBy, Value: sortOrder}}}})
	}

	// Pagination
	if options.Page > 0 && options.PageSize > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (options.Page - 1) * options.PageSize}})
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: options.PageSize}})
	}

	// Remove reviews array from final result
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.D{
		{Key: "reviews", Value: 0},
	}}})

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var businesses []Business
	if err = cursor.All(ctx, &businesses); err != nil {
		return nil, 0, err
	}

	return businesses, totalCount, nil
}

func (s *mongoStore) GetBusinessByID(id string) (*Business, error) {
	idPrim, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var business Business
	filter := bson.M{"_id": idPrim}
	err = s.coll.FindOne(context.TODO(), filter).Decode(&business)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &business, nil
}

// Inserts the given business into the collection using the given ctx.
func (s *mongoStore) CreateBusiness(ctx context.Context, b *Business) error {
	b.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	b.UpdatedAt = b.CreatedAt

	_, err := s.coll.InsertOne(ctx, b)
	return err
}
