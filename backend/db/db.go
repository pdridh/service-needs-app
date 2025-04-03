package db

import (
	"context"
	"log"
	"time"

	"github.com/pdridh/service-needs-app/backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func DisconnectFromDB() {
	if client == nil {
		log.Fatal("Accessing client before connection")
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal("disconnect error", err)
	}
	log.Print("Disconnect from the DB")
}

func ensureGeospatialIndex(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get existing indexes
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return err
	}

	// Check if our index already exists
	var indexResults []bson.M
	if err := cursor.All(ctx, &indexResults); err != nil {
		return err
	}

	for _, index := range indexResults {
		if name, ok := index["name"].(string); ok && name == "location_2dsphere" {
			return nil
		}
	}

	// Create index only if it doesn't exist
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "location", Value: "2dsphere"}},
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

func ConnectToDB() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Server().DatabaseURI))
	if err != nil {
		log.Println("error connecting to the db")
		log.Fatal(err)
	}

	if err := ensureGeospatialIndex(GetCollectionFromDB(config.Server().DatabaseName, config.Server().BusinessCollectionName)); err != nil {
		DisconnectFromDB()
		log.Fatal(err)
	}

	log.Print("Succesfully connected to the DB")
}

func GetClient() *mongo.Client {
	if client == nil {
		log.Fatal("Accessing client before connection")
	}

	return client
}

// Helper to access a collection from a db using the client
func GetCollectionFromDB(db string, coll string) *mongo.Collection {
	return GetClient().Database(db).Collection(coll)
}
