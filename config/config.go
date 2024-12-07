package config

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once      sync.Once
	uri       string
	dbName    string
	jwtSecret string
)

type Config struct {
	MongoClient *mongo.Client
	Database    *mongo.Database
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logrus.Printf("Warning: .env file not found: %v", err)
	}
	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASSWORD")
	dbName = os.Getenv("MONGO_DATABASE")
	jwtSecret = os.Getenv("JWT_SECRET")

	if mongoUser == "" || mongoPass == "" || dbName == "" || jwtSecret == "" {
		logrus.Fatal("MONGO_USER, MONGO_PASSWORD, MONGO_DATABASE and JWT_SECRET are required")
	}

	host := "localhost"
	if os.Getenv("DOCKER_CONTAINER") == "true" {
		host = "mongo"
	}

	uri = fmt.Sprintf("mongodb://%s:%s@%s:27017/%s?authSource=admin",
		mongoUser,
		mongoPass,
		host,
		dbName,
	)
}

func ConnectToMongo() (*Config, error) {
	LoadEnv()
	logrus.Printf("Connect to monga: uri: %s\n", uri)
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(dbName)
	return &Config{
		MongoClient: client,
		Database:    database,
	}, nil
}

func (c *Config) CloseMongo() {
	if c.MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.MongoClient.Disconnect(ctx); err != nil {
			logrus.Printf("Error while disconnecting from MongoDB: %v", err)
		}
	}
}

func GetJWTSecret() string {
	once.Do(LoadEnv)
	return jwtSecret
}
