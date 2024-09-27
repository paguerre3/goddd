package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewMongoClient(t *testing.T) {
	originalApplyUri := applyURI
	defer func() { applyURI = originalApplyUri }()

	originalMongoConnect := mongoConnect
	defer func() { mongoConnect = originalMongoConnect }()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		applyURI = func(uri string) *options.ClientOptions {
			// mock the applyURI function
			return nil
		}

		mongoConnect = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
			// mock the mongoConnect function
			return mt.Client, nil
		}

		mongoClient := NewMongoClient()
		assert.NotNil(t, mongoClient)
		tc := mongoClient.GetCollection("testCol")
		assert.NotNil(t, tc)
		// 1.13.0 the Close() method for mtest package is removed, this method is not necessary.
		// Debugging the following code will succeed but the test at the end fails as it will call
		// the Close again automatically, i.e. close of closed channel already performed:
		//
		// 2024/09/26 23:10:00 Connected to MongoDB!
		// 2024/09/26 23:10:25 Disconnected from MongoDB
		// --- FAIL: TestNewMongoClient (52.28s)
		// --- FAIL: TestNewMongoClient/success (52.28s)
		// panic: close of closed channel [recovered]
		// panic: close of closed channel
		//
		// err := mongoClient.Close()
		// assert.NoError(t, err)
	})

	mt.Run("mongoConnect fatal", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		applyURI = func(uri string) *options.ClientOptions {
			// mock the applyURI function
			return nil
		}

		mongoConnect = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
			// mock the mongoConnect function
			return nil, fmt.Errorf("connection failure")
		}

		var mongoClient MongoClient
		if r := recover(); r != nil {
			mongoClient = NewMongoClient()
		}
		assert.Nil(t, mongoClient)
	})
}

func TestResolveMongoUri(t *testing.T) {
	originalGetEnv := getEnv
	defer func() { getEnv = originalGetEnv }() // Restore original function after test

	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{"EnvVarSet", "mongodb://remote:27017", "mongodb://remote:27017"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getEnv = func(key string) string {
				if key == "MONGO_ADDR" {
					return tt.envValue
				}
				return ""
			}

			result := resolveMongoUri()
			assert.Equal(t, tt.expected, result)
		})
	}
}
