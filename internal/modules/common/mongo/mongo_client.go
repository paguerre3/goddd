package mongo

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient interface {
	GetCollection(collectionName string) *mongo.Collection
	Close() error
}

type mongoClient struct {
	client   *mongo.Client
	database *mongo.Database
}

var (
	mongoInstance *mongoClient
	once          sync.Once
	applyURI      = options.Client().ApplyURI
	mongoConnect  = mongo.Connect
	getEnv        = os.Getenv
)

const (
	timeout       = 10 * time.Second
	mongoLocalUri = "mongodb://root:pass@localhost:27017"
	dbName        = "padeldb"
)

// NewMongoClient initializes and returns a singleton MongoClient
func NewMongoClient() MongoClient {
	once.Do(func() {
		mongoURI := resolveMongoUri()

		clientOptions := applyURI(mongoURI)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		client, err := mongoConnect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Check the MongoDB connection
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatalf("Could not ping to MongoDB: %v", err)
		}

		log.Println("Connected to MongoDB!")

		database := client.Database(dbName) // Database name

		mongoInstance = &mongoClient{
			client:   client,
			database: database,
		}
	})

	return mongoInstance
}

// GetCollection returns a Mongo collection by name
func (m *mongoClient) GetCollection(collectionName string) *mongo.Collection {
	return m.database.Collection(collectionName)
}

// Close gracefully disconnects the MongoDB client
func (m *mongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		return err
	}

	log.Println("Disconnected from MongoDB")
	return nil
}

func resolveMongoUri() string {
	uri := getEnv("MONGO_ADDR")
	if len(uri) > 0 {
		return uri
	}
	return mongoLocalUri
}
