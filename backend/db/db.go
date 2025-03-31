package db

import (
	"context"
	"log"

	"github.com/pdridh/service-needs-app/backend/config"
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

func ConnectToDB() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Server().DatabaseURI))
	if err != nil {
		log.Println("error connecting to the db")
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
