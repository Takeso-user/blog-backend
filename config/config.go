package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once       sync.Once
	uri        string
	dbName     string
	jwtSecret  string
	mongoUser  string
	mongoPass  string
	secretsDir = "/run/secrets/"
)

type Config struct {
	MongoClient *mongo.Client
	Database    *mongo.Database
}

func LoadSecrets() {
	log.Println("Loading secrets from Docker Secrets directory...")
	mongoUser = readSecret("mongo_user")
	mongoPass = readSecret("mongo_password")
	jwtSecret = readSecret("jwt_secret")

	if mongoUser == "" || mongoPass == "" {
		log.Println("Secrets not found. Falling back to .env file.")
		LoadEnv()
	} else {
		log.Println("Secrets loaded successfully.")
		uri = fmt.Sprintf("mongodb://%s:%s@mongo:27017", mongoUser, mongoPass)
		dbName = os.Getenv("MONGO_DATABASE") // Оставляем из окружения или .env
	}
}

func readSecret(name string) string {
	filePath := secretsDir + name
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Secret %s not found: %v", name, err)
		return ""
	}
	return string(data)
}

func LoadEnv() {
	log.Println("Loading environment variables from .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Environment variables loaded successfully.")

	uri = os.Getenv("MONGO_URI")
	dbName = os.Getenv("MONGO_DATABASE")
}

func ConnectToMongo() (*Config, error) {
	once.Do(func() {
		LoadSecrets()
	})

	if uri == "" || dbName == "" {
		return nil, fmt.Errorf("MONGO_URI or MONGO_DATABASE is not set")
	}

	log.Println("Creating MongoDB client options...")
	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)

	log.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	log.Println("Pinging MongoDB...")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully.")
	database := client.Database(dbName)

	return &Config{
		MongoClient: client,
		Database:    database,
	}, nil
}

func (c *Config) CloseMongo() {
	log.Println("Disconnecting from MongoDB...")
	if err := c.MongoClient.Disconnect(context.TODO()); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB successfully.")
	}
}
