package mongo

import (
	"fmt"
	"testing"

	common "github.com/paguerre3/goddd/internal/modules/common/mongo"
	"github.com/paguerre3/goddd/internal/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

const (
	testDbName          = "testdb"
	testPlayersNs       = testDbName + "." + playersColName
	testPlayerCouplesNs = testDbName + "." + playerCouplesColName
)

type mongoClientMock struct {
	client   *mongo.Client
	database *mongo.Database
}

func (m *mongoClientMock) GetCollection(collectionName string) *mongo.Collection {
	return m.database.Collection(collectionName)
}

func (m *mongoClientMock) Close() error {
	// 1.13.0 the Close() method for mtest package is removed, this method is not necessary
	return nil
}

func NewMongoClientMock(client *mongo.Client) common.MongoClient {
	return &mongoClientMock{
		client:   client,
		database: client.Database(testDbName),
	}
}

func TestMongoPlayerRepository_Save_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// 1.13.0 the Close() method for mtest package is removed, this method is not necessary

	mt.Run("Save player successfully", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)
		player := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Save(player)
		assert.NoError(t, err, "Expected no error when saving player")

	})
}

func TestMongoPlayerRepository_Save_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to save player", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)
		player := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		err := repo.Save(player)
		assert.Error(t, err, "Expected error when saving player")
	})
}

func TestMongoPlayerRepository_FindByID_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find player by ID successfully", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)
		player := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "id", Value: player.ID},
			{Key: "firstName", Value: player.FirstName},
			{Key: "lastName", Value: player.LastName},
			{Key: "email", Value: player.Email},
		}))

		result, err := repo.FindByID("1")
		assert.NoError(t, err, "Expected no error when finding player by ID")
		assert.Equal(t, player, result, "Expected player to match")
	})
}

func TestMongoPlayerRepository_FindByID_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to find player by ID", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, testPlayersNs, mtest.FirstBatch))

		result, err := repo.FindByID("1")
		assert.Error(t, err, "Expected error when finding player by ID")
		assert.Equal(t, domain.Player{}, result, "Expected result to be empty player")
	})
}

func TestMongoPlayerRepository_FindByEmail_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find player by email successfully", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)
		player := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "id", Value: player.ID},
			{Key: "firstName", Value: player.FirstName},
			{Key: "lastName", Value: player.LastName},
			{Key: "email", Value: player.Email},
		}))

		result, err := repo.FindByEmail("john.doe@example.com")
		assert.NoError(t, err, "Expected no error when finding player by email")
		assert.Equal(t, player, result, "Expected player to match")
	})
}

func TestMongoPlayerRepository_FindByEmail_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to find player by email", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, testPlayersNs, mtest.FirstBatch))

		result, err := repo.FindByEmail("john.doe@example.com")
		assert.Error(t, err, "Expected error when finding player by email")
		assert.Equal(t, domain.Player{}, result, "Expected result to be empty player")
	})
}

func TestMongoPlayerRepository_FindByLastName_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find players by last name successfully", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)
		player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
		player2 := domain.Player{ID: "2", FirstName: "Juan", LastName: "Doe", Email: "juan.doe@example.com"}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "id", Value: player1.ID},
			{Key: "firstName", Value: player1.FirstName},
			{Key: "lastName", Value: player1.LastName},
			{Key: "email", Value: player1.Email},
		}), mtest.CreateCursorResponse(1, testPlayersNs, mtest.NextBatch, bson.D{
			{Key: "id", Value: player2.ID},
			{Key: "firstName", Value: player2.FirstName},
			{Key: "lastName", Value: player2.LastName},
			{Key: "email", Value: player2.Email},
		}))

		result, err := repo.FindByLastName("Doe")
		assert.NoError(t, err, "Expected no error when finding players by last name")
		assert.Equal(t, []domain.Player{player1, player2}, result, "Expected players to match")
	})
}

func TestMongoPlayerRepository_FindByLastName_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to find players by last name", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(-1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "id", Value: "1"},
			{Key: "lastName", Value: "Smith"},
			{Key: "age", Value: "invalidAge"},
		}))

		result, err := repo.FindByLastName("Smith")
		assert.Error(t, err, "Expected error when finding players by last name")
		assert.Nil(t, result, "Expected result to be nil")
	})
}

func TestMongoPlayerRepository_Update_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Update player successfully", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)
		player := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Update(player)
		assert.NoError(t, err, "Expected no error when updating player")
	})
}

func TestMongoPlayerRepository_Update_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to update player", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)
		player := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		err := repo.Update(player)
		assert.Error(t, err, "Expected error when updating player")
	})
}

func TestMongoPlayerRepository_Delete_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Delete player successfully", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Delete("1")
		assert.NoError(t, err, "Expected no error when deleting player")
	})
}

func TestMongoPlayerRepository_Delete_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to delete player", func(mt *mtest.T) {
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(mongoClientMock)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "delete error",
		}))

		err := repo.Delete("1")
		assert.Error(t, err, "Expected error when deleting player")
	})
}

func TestMongoPlayerCoupleRepository_Save(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		playerCouple := domain.PlayerCouple{ID: "c1", Player1: player1, Player2: player2}
		err := repo.Save(playerCouple)
		assert.NoError(t, err, "Expected no error when saving player couple")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		playerCouple := domain.PlayerCouple{ID: "c1", Player1: player1, Player2: player2}
		err := repo.Save(playerCouple)
		assert.Error(t, err, "Expected error when saving player couple")
	})
}

func TestMongoPlayerCoupleRepository_FindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}
	playerCouple := domain.PlayerCouple{ID: "c1", Player1: player1, Player2: player2}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayerCouplesNs, mtest.FirstBatch, bson.D{
			{Key: "id", Value: playerCouple.ID},
			{Key: "player1", Value: bson.D{
				{Key: "id", Value: player1.ID},
				{Key: "firstName", Value: player1.FirstName},
				{Key: "lastName", Value: player1.LastName},
				{Key: "email", Value: player1.Email},
			}},
			{Key: "player2", Value: bson.D{
				{Key: "id", Value: player2.ID},
				{Key: "firstName", Value: player2.FirstName},
				{Key: "lastName", Value: player2.LastName},
				{Key: "email", Value: player2.Email},
			}},
		}))

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		result, err := repo.FindByID(playerCouple.ID)
		assert.NoError(t, err, "Expected no error when finding player couple by ID")
		assert.Equal(t, playerCouple, result, "Expected player couple to match")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, testPlayerCouplesNs, mtest.FirstBatch))

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		_, err := repo.FindByID("c2")
		assert.Error(t, err, "Expected error when finding by player couple ID")
	})
}

func TestMongoPlayerCoupleRepository_FindByPrefixes(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}
	cid := fmt.Sprintf("%s-%s-coupleMockId", player1.LastName, player2.LastName)
	playerCouple := domain.PlayerCouple{ID: cid, Player1: player1, Player2: player2}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayerCouplesNs, mtest.FirstBatch, bson.D{
			{Key: "id", Value: playerCouple.ID},
			{Key: "player1", Value: bson.D{
				{Key: "id", Value: player1.ID},
				{Key: "firstName", Value: player1.FirstName},
				{Key: "lastName", Value: player1.LastName},
				{Key: "email", Value: player1.Email},
			}},
			{Key: "player2", Value: bson.D{
				{Key: "id", Value: player2.ID},
				{Key: "firstName", Value: player2.FirstName},
				{Key: "lastName", Value: player2.LastName},
				{Key: "email", Value: player2.Email},
			}},
		}))

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		result, err := repo.FindByPrefixes(player1.LastName, player2.LastName)
		assert.NoError(t, err, "Expected no error when finding player couple by prefixes")
		assert.Equal(t, []domain.PlayerCouple{playerCouple}, result, "Expected player couples to match")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayerCouplesNs, mtest.FirstBatch, bson.D{
			{Key: "id", Value: playerCouple.ID},
			{Key: "player1", Value: bson.D{
				{Key: "id", Value: player1.ID},
				{Key: "firstName", Value: player1.FirstName},
				{Key: "lastName", Value: player1.LastName},
				{Key: "email", Value: player1.Email},
			}},
			{Key: "player2", Value: bson.D{
				{Key: "id", Value: player2.ID},
				{Key: "firstName", Value: player2.FirstName},
				{Key: "lastName", Value: player2.LastName},
				{Key: "email", Value: player2.Email},
				// Decode error:
				{Key: "age", Value: "invalidAge"},
			}},
		}))

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		result, err := repo.FindByPrefixes(player1.LastName, player2.LastName)
		assert.Error(t, err, "Expected error when finding player couple by prefixes")
		assert.Nil(t, result, "Expected result to be nil")
	})
}

func TestMongoPlayerCoupleRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}
	playerCouple := domain.PlayerCouple{ID: "c1", Player1: player1, Player2: player2}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)

		err := repo.Update(playerCouple)
		assert.NoError(t, err, "Expected no error when updating player couple")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		err := repo.Update(playerCouple)
		assert.Error(t, err, "Expected error when updating player couple")
	})
}

func TestMongoPlayerCoupleRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		err := repo.Delete("1")
		assert.NoError(t, err, "Expected no error when deleting player couple")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "delete error",
		}))

		mongoClientMock := NewMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(mongoClientMock)
		err := repo.Delete("1")
		assert.Error(t, err, "Expected error when deleting player couple")
	})
}
