package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoInstance *MongoDB // Exported MongoDB instance

// MongoDB holds the database connection
type MongoDB struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

// InitMongoDB initializes the MongoDB connection
func InitMongoDB(uri, dbName, collectionName string) error {
	var err error
	MongoInstance, err = NewMongoDB(uri, dbName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to initialize MongoDB: %v", err)
	}
	return nil
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri, dbName, collectionName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &MongoDB{
		Client:     client,
		Database:   database,
		Collection: collection,
	}, nil
}
