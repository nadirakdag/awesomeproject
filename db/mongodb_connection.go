package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Db *mongo.Database
}

func NewMongoConnection(serverURL, databaseName string) (*MongoConnection, error) {

	opts := options.Client()
	opts.ApplyURI(serverURL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), opts)

	if err != nil {
		return nil, err
	}

	return &MongoConnection{Db: client.Database(databaseName)}, nil
}
