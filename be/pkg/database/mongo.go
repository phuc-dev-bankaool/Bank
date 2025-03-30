package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MongoDB client instance
var Client *mongo.Client

// ConnectMongo establishes a connection to MongoDB
func ConnectMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions) // Đã thêm ctx vào hàm Connect
	if err != nil {
		return nil, err
	}

	// Kiểm tra kết nối
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	Client = client
	return client, nil
}

// GetCollection returns a MongoDB collection
func GetCollection(dbName, collName string) *mongo.Collection {
	return Client.Database(dbName).Collection(collName)
}
