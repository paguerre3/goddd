package main

import (
	"context"
	"log"

	infrastructure "github.com/paguerre3/goddd/modules/player-couple/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "padledb"
	// TODO: get from environment variables (set in dockerfile):
	dbUri = "mongodb://localhost:27017"
)

func main() {
	clientOptions := options.Client().ApplyURI(dbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	playerRepo := infrastructure.NewMongoPlayerRepository(client, dbName, "player")
	playerRepo.FindByEmail("agus.tapia@gmail.com")
	playerCoupleRepo := infrastructure.NewMongoPlayerCoupleRepository(client, dbName, "player_couple")
	playerCoupleRepo.FindByPrefixes("Tapia", "")
}
