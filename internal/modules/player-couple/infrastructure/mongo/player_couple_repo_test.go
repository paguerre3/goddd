package mongo

import (
	"fmt"
	"testing"

	common "github.com/paguerre3/goddd/internal/modules/common/mongo"
	"github.com/paguerre3/goddd/internal/modules/common/utils"
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
	mockId              = "mock-id"
)

type idGenMock struct {
}

func (i *idGenMock) GenerateID() string {
	return mockId
}

func (i *idGenMock) GenerateIDWithPrefixes(prefix1 string, prefix2 string) string {
	return fmt.Sprintf("%s-%s-%s", prefix1, prefix2, i.GenerateID())
}

func newIdGenMock() utils.IDGenerator {
	return &idGenMock{}
}

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

func newMongoClientMock(client *mongo.Client) common.MongoClient {
	return &mongoClientMock{
		client:   client,
		database: client.Database(testDbName),
	}
}

func TestMongoPlayerRepository_Save_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// 1.13.0 the Close() method for mtest package is removed, this method is not necessary

	mt.Run("Save player successfully", func(mt *mtest.T) {
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)
		player, err := domain.NewPlayer("john.doe@example.com", nil, "John", "Doe", nil)
		assert.NoError(t, err)
		assert.Equal(t, "", player.ID)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err = repo.Upsert(player)
		// generated ID set in repository implies a Save():
		assert.Equal(t, mockId, player.ID)
		assert.NoError(t, err, "Expected no error when saving player")

	})
}

func TestMongoPlayerRepository_Save_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to save player", func(mt *mtest.T) {
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)
		player, err := domain.NewPlayer("john.doe@example.com", nil, "John", "Doe", nil)
		assert.NoError(t, err)
		// generated ID set in repository implies a Save():
		assert.Equal(t, "", player.ID)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		err = repo.Upsert(player)
		assert.Error(t, err, "Expected error when saving player")
		assert.Equal(t, "", player.ID)
	})
}

func TestMongoPlayerRepository_FindByID_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find player by ID successfully", func(mt *mtest.T) {
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)
		player, err := domain.NewPlayer("john.doe@example.com", nil, "John", "Doe", nil)
		assert.NoError(t, err)
		assert.Equal(t, "", player.ID)
		// Mock a player with a generated ID set in the repository:
		player.ID = idGen.GenerateID()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "_id", Value: player.ID},
			{Key: "firstName", Value: player.FirstName},
			{Key: "lastName", Value: player.LastName},
			{Key: "email", Value: player.Email},
		}))

		result, err := repo.FindByID(mockId)
		assert.NoError(t, err, "Expected no error when finding player by ID")
		assert.Equal(t, *player, result, "Expected player to match")
	})
}

func TestMongoPlayerRepository_FindByID_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find player by ID not found", func(mt *mtest.T) {
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(newIdGenMock(), mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, testPlayersNs, mtest.FirstBatch))

		result, err := repo.FindByID(mockId)
		assert.NoError(t, err)
		assert.Equal(t, domain.Player{}, result, "Expected result to be empty player")
	})
}

func TestMongoPlayerRepository_FindByID_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to find player by ID", func(mt *mtest.T) {
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(newIdGenMock(), mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(-1, testPlayersNs, mtest.FirstBatch))

		result, err := repo.FindByID(mockId)
		assert.Error(t, err, "Expected error when finding player by ID")
		assert.Equal(t, domain.Player{}, result, "Expected result to be empty player")
	})
}

func TestMongoPlayerRepository_FindByEmail_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find player by email successfully", func(mt *mtest.T) {
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)
		player := domain.Player{ID: idGen.GenerateID(), FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "_id", Value: player.ID},
			{Key: "firstName", Value: player.FirstName},
			{Key: "lastName", Value: player.LastName},
			{Key: "email", Value: player.Email},
		}))

		result, err := repo.FindByEmail("john.doe@example.com")
		assert.NoError(t, err, "Expected no error when finding player by email")
		assert.Equal(t, player, result, "Expected player to match")
	})
}

func TestMongoPlayerRepository_FindByEmail_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find player by email not found", func(mt *mtest.T) {
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(newIdGenMock(), mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, testPlayersNs, mtest.FirstBatch))

		result, err := repo.FindByEmail("john.doe@example.com")
		assert.NoError(t, err)
		assert.Equal(t, domain.Player{}, result, "Expected result to be empty player")
	})
}

func TestMongoPlayerRepository_FindByEmail_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to find player by email", func(mt *mtest.T) {
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(newIdGenMock(), mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(-1, testPlayersNs, mtest.FirstBatch))

		result, err := repo.FindByEmail("john.doe@example.com")
		assert.Error(t, err, "Expected error when finding player by email")
		assert.Equal(t, domain.Player{}, result, "Expected result to be empty player")
	})
}

func TestMongoPlayerRepository_FindByLastName_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Find players by last name successfully", func(mt *mtest.T) {
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)
		player1 := domain.Player{ID: idGen.GenerateID() + "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
		player2 := domain.Player{ID: idGen.GenerateID() + "2", FirstName: "Juan", LastName: "Doe", Email: "juan.doe@example.com"}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "_id", Value: player1.ID},
			{Key: "firstName", Value: player1.FirstName},
			{Key: "lastName", Value: player1.LastName},
			{Key: "email", Value: player1.Email},
		}), mtest.CreateCursorResponse(1, testPlayersNs, mtest.NextBatch, bson.D{
			{Key: "_id", Value: player2.ID},
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
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)

		mt.AddMockResponses(mtest.CreateCursorResponse(-1, testPlayersNs, mtest.FirstBatch, bson.D{
			{Key: "_id", Value: idGen.GenerateID()},
			{Key: "lastName", Value: "Smith"},
			{Key: "age", Value: "invalidAgeDecode"},
		}))

		result, err := repo.FindByLastName("Smith")
		assert.Error(t, err, "Expected error when finding players by last name")
		assert.Nil(t, result, "Expected result to be nil")
	})
}

func TestMongoPlayerRepository_Upsert_Update_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Update player successfully", func(mt *mtest.T) {
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)
		excpectedId := idGen.GenerateID()
		player := domain.Player{ID: excpectedId, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Update inplies ID already set previous to the Upsert method call.
		err := repo.Upsert(&player)
		// NOT new autogenerated ID set in repository implies an Update():
		assert.Equal(t, excpectedId, player.ID)
		assert.NoError(t, err, "Expected no error when updating player")
	})
}

func TestMongoPlayerRepository_Upsert_Update_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to update player", func(mt *mtest.T) {
		idGen := newIdGenMock()
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(idGen, mongoClientMock)
		// Update inplies ID already set previous to the Upsert method call.
		player := domain.Player{ID: idGen.GenerateID(), FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		err := repo.Upsert(&player)
		assert.Error(t, err, "Expected error when updating player")
	})
}

func TestMongoPlayerRepository_Delete_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Delete player successfully", func(mt *mtest.T) {
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(newIdGenMock(), mongoClientMock)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Delete("1")
		assert.NoError(t, err, "Expected no error when deleting player")
	})
}

func TestMongoPlayerRepository_Delete_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fail to delete player", func(mt *mtest.T) {
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerRepository(newIdGenMock(), mongoClientMock)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "delete error",
		}))

		err := repo.Delete("1")
		assert.Error(t, err, "Expected error when deleting player")
	})
}

func TestMongoPlayerCoupleRepository_Upsert_Save(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	idGen := newIdGenMock()
	genPlayerId1 := idGen.GenerateID() + "1"
	genPlayerId2 := idGen.GenerateID() + "2"
	player1 := domain.Player{ID: genPlayerId1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: genPlayerId2, FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(idGen, mongoClientMock)
		// couple ID not set in domain implies a Save() before Upsert() call:
		playerCouple, err := domain.NewPlayerCouple(player1, player2, nil)
		assert.Equal(t, "", playerCouple.ID)
		err = repo.Upsert(playerCouple)
		assert.NoError(t, err, "Expected no error when saving player couple")
		// couple ID set in repository implies a Save() inside Upsert() call:
		expectedCoupleID := idGen.GenerateIDWithPrefixes(player1.LastName, player2.LastName)
		assert.Equal(t, expectedCoupleID, playerCouple.ID)
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(idGen, mongoClientMock)
		// couple ID not set in domain implies a Save() before Upsert() call:
		playerCouple, err := domain.NewPlayerCouple(player1, player2, nil)
		assert.Equal(t, "", playerCouple.ID)

		err = repo.Upsert(playerCouple)
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
			{Key: "_id", Value: playerCouple.ID},
			{Key: "player1", Value: bson.D{
				{Key: "_id", Value: player1.ID},
				{Key: "firstName", Value: player1.FirstName},
				{Key: "lastName", Value: player1.LastName},
				{Key: "email", Value: player1.Email},
			}},
			{Key: "player2", Value: bson.D{
				{Key: "_id", Value: player2.ID},
				{Key: "firstName", Value: player2.FirstName},
				{Key: "lastName", Value: player2.LastName},
				{Key: "email", Value: player2.Email},
			}},
		}))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(newIdGenMock(), mongoClientMock)
		result, err := repo.FindByID(playerCouple.ID)
		assert.NoError(t, err, "Expected no error when finding player couple by ID")
		assert.Equal(t, playerCouple, result, "Expected player couple to match")
	})

	mt.Run("not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, testPlayerCouplesNs, mtest.FirstBatch))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(newIdGenMock(), mongoClientMock)
		pc, err := repo.FindByID("c2")
		assert.NoError(t, err, "Expected no error when player couple not found")
		assert.Equal(t, domain.PlayerCouple{}, pc, "Expected empty player couple")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(-1, testPlayerCouplesNs, mtest.FirstBatch))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(newIdGenMock(), mongoClientMock)
		_, err := repo.FindByID("c2")
		assert.Error(t, err, "Expected error when finding by player couple ID")
	})
}

func TestMongoPlayerCoupleRepository_FindByPrefixes(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	idGen := newIdGenMock()
	player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}
	cid := idGen.GenerateIDWithPrefixes(player1.LastName, player2.LastName)
	playerCouple := domain.PlayerCouple{ID: cid, Player1: player1, Player2: player2}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayerCouplesNs, mtest.FirstBatch, bson.D{
			{Key: "_id", Value: playerCouple.ID},
			{Key: "player1", Value: bson.D{
				{Key: "_id", Value: player1.ID},
				{Key: "firstName", Value: player1.FirstName},
				{Key: "lastName", Value: player1.LastName},
				{Key: "email", Value: player1.Email},
			}},
			{Key: "player2", Value: bson.D{
				{Key: "_id", Value: player2.ID},
				{Key: "firstName", Value: player2.FirstName},
				{Key: "lastName", Value: player2.LastName},
				{Key: "email", Value: player2.Email},
			}},
		}))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(idGen, mongoClientMock)
		result, err := repo.FindByPrefixes(player1.LastName, player2.LastName)
		assert.NoError(t, err, "Expected no error when finding player couple by prefixes")
		assert.Equal(t, []domain.PlayerCouple{playerCouple}, result, "Expected player couples to match")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(1, testPlayerCouplesNs, mtest.FirstBatch, bson.D{
			{Key: "_id", Value: playerCouple.ID},
			{Key: "player1", Value: bson.D{
				{Key: "_id", Value: player1.ID},
				{Key: "firstName", Value: player1.FirstName},
				{Key: "lastName", Value: player1.LastName},
				{Key: "email", Value: player1.Email},
			}},
			{Key: "player2", Value: bson.D{
				{Key: "_id", Value: player2.ID},
				{Key: "firstName", Value: player2.FirstName},
				{Key: "lastName", Value: player2.LastName},
				{Key: "email", Value: player2.Email},
				// Decode error:
				{Key: "age", Value: "invalidAgeDecode"},
			}},
		}))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(idGen, mongoClientMock)
		result, err := repo.FindByPrefixes(player1.LastName, player2.LastName)
		assert.Error(t, err, "Expected error when finding player couple by prefixes")
		assert.Nil(t, result, "Expected result to be nil")
	})
}

func TestMongoPlayerCoupleRepository_Upsert_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	idGen := newIdGenMock()
	player1 := domain.Player{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	player2 := domain.Player{ID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"}
	coupleIdExpected := idGen.GenerateIDWithPrefixes(player1.LastName, player2.LastName)
	// couple ID set in domain previous to the Upsert() call iplies an Update():
	playerCouple := domain.PlayerCouple{ID: coupleIdExpected, Player1: player1, Player2: player2}

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(idGen, mongoClientMock)

		assert.Equal(t, coupleIdExpected, playerCouple.ID)
		err := repo.Upsert(&playerCouple)
		assert.NoError(t, err, "Expected no error when updating player couple")
		assert.Equal(t, coupleIdExpected, playerCouple.ID)
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "duplicate key error",
		}))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(idGen, mongoClientMock)

		assert.Equal(t, coupleIdExpected, playerCouple.ID)
		err := repo.Upsert(&playerCouple)
		assert.Error(t, err, "Expected error when updating player couple")
	})
}

func TestMongoPlayerCoupleRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(newIdGenMock(), mongoClientMock)
		err := repo.Delete("1")
		assert.NoError(t, err, "Expected no error when deleting player couple")
	})

	mt.Run("failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "delete error",
		}))

		mongoClientMock := newMongoClientMock(mt.Client)
		repo := NewMongoPlayerCoupleRepository(newIdGenMock(), mongoClientMock)
		err := repo.Delete("1")
		assert.Error(t, err, "Expected error when deleting player couple")
	})
}
