package config

import (
	"backend/constants"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func ConnectDB() {
	mongoURI := GetEnv("MONGO_URI")
	dbName := constants.DB_NAME
	if mongoURI == "" || dbName == "" {
		log.Fatal("MONGO_URI or DB_NAME is not defined in environment variables")
	}
	// MongoDB connection options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(context.Background())
	dbClient = client

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to mongoDB database...")
}

// GetCollection retrieves a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	if dbClient == nil {
		log.Fatal("Database connection is not initialized")
	}
	return dbClient.Database(constants.DB_NAME).Collection(collectionName)
}
