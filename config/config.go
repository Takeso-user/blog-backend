package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once   sync.Once
	uri    string
	dbName string
)

type Config struct {
	MongoClient *mongo.Client
	Database    *mongo.Database
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func ConnectToMongo() (*Config, error) {
	once.Do(
		func() {
			uri = os.Getenv("MONGO_URI")
			dbName = os.Getenv("MONGO_DATABASE")
		})

	if uri == "" || dbName == "" {
		return nil, fmt.Errorf("MONGO_URI or MONGO_DATABASE is not set in the environment")
	}

	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB")

	database := client.Database(dbName)

	return &Config{
		MongoClient: client,
		Database:    database,
	}, nil
}

func (c *Config) CloseMongo() {
	if err := c.MongoClient.Disconnect(context.TODO()); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB")
	}
}
