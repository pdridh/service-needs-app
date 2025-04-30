package user

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	GetUsers(ctx context.Context, options QueryOptions) ([]User, int64, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(ctx context.Context, u *User) error
}

type mongoStore struct {
	coll *mongo.Collection
}

// Creates a new mongo store with the collection and returns a ptr to it.
func NewMongoStore(coll *mongo.Collection) *mongoStore {
	return &mongoStore{coll: coll}
}

func (s *mongoStore) GetUsers(ctx context.Context, options QueryOptions) ([]User, int64, error) {
	// This shit is so horrible i cannot even start, making me regret using mongo ong
	pipeline := mongo.Pipeline{}

	if options.Search != "" {
		searchFields := []string{"email"}
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
		log.Println("Pagination should work idk")
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (options.Page - 1) * options.PageSize}})
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: options.PageSize}})
	}

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// Helper function to get a user by a dynamic field and its value
// Returns a ptr to a User struct or nil if not found
func (s *mongoStore) getUserByField(field string, value string) (*User, error) {
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
func (s *mongoStore) GetUserByEmail(email string) (*User, error) {
	return s.getUserByField("email", email)
}

// Given a user ptr adds required fields and inserts into the coll using the given context.
// DOESNT DO ANY CHECKS.
// returns error if any or nil.
func (s *mongoStore) CreateUser(ctx context.Context, u *User) error {

	u.ID = primitive.NewObjectID()
	u.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	u.UpdatedAt = u.CreatedAt

	_, err := s.coll.InsertOne(ctx, u)
	return err
}
