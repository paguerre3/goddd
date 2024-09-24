package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/benweissmann/memongo"
	"github.com/paguerre3/goddd/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoPlayerCoupleRepository(t *testing.T) {
	// Start an in-memory MongoDB server
	// https://www.mongodb.com/resources/products/mongodb-version-history
	mongoServer, err := memongo.Start("7.0.0")
	if err != nil {
		t.Fatalf("failed to start in-memory MongoDB server: %v", err)
	}
	defer mongoServer.Stop()

	// Connect to the in-memory MongoDB server
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoServer.URI()))
	if err != nil {
		t.Fatalf("failed to connect to in-memory MongoDB server: %v", err)
	}

	repo := NewMongoPlayerCoupleRepository(client, "memdb", "player_couples")

	player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}
	cid := fmt.Sprintf("%s-%s-coupleMockId", player1.LastName, player2.LastName)
	playerCouple := domain.PlayerCouple{ID: cid, Player1: player1, Player2: player2}

	t.Run("Save", func(t *testing.T) {
		err := repo.Save(playerCouple)
		assert.NoError(t, err, "Expected no error when saving player couple in memory")
	})

	t.Run("FindByID", func(t *testing.T) {
		result, err := repo.FindByID(cid)
		assert.NoError(t, err, "Expected no error when finding player couple by ID in memory")
		assert.Equal(t, playerCouple, result, "Expected player couple to match")
	})
}
