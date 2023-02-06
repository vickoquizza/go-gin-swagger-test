package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoCollection(URI string) (*mongo.Collection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))

	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	err = client.Connect(ctx)

	if err != nil {
		return nil, err
	}

	collection := client.Database("test").Collection("accounts")

	return collection, nil
}
