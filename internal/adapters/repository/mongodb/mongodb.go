package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collections names
const (
	ProductCollection = "products"
	OrderCollection   = "orders"
	UserCollection    = "users"
	AddressCollection = "addresses"
	SettingCollection = "settings"
)

// NewMongoDBConnection establishes a connection to MongoDB
func NewMongoDBConnection(uri, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(database), nil
}
