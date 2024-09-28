package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	common "github.com/paguerre3/goddd/internal/modules/common/mongo"
	"github.com/paguerre3/goddd/internal/modules/common/utils"
	"github.com/paguerre3/goddd/internal/modules/player-couple/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	timeout              = 5 * time.Second
	playersColName       = "players"
	playerCouplesColName = "player_couples"
)

type mongoPlayerRepository struct {
	idGen      utils.IDGenerator
	collection *mongo.Collection
}

type mongoPlayerCoupleRepository struct {
	idGen      utils.IDGenerator
	collection *mongo.Collection
}

func NewMongoPlayerRepository(idGen utils.IDGenerator, client common.MongoClient) domain.PlayerRepository {
	collection := client.GetCollection(playersColName)
	return &mongoPlayerRepository{
		idGen:      idGen,
		collection: collection,
	}
}

func (r *mongoPlayerRepository) Upsert(player *domain.Player) error {
	if player == nil {
		return errors.New("player is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// DDD repository principle.
	if len(player.ID) > 0 {
		_, err := r.collection.UpdateOne(ctx, bson.M{"_id": player.ID}, bson.M{"$set": player})
		return err
	}
	player.ID = r.idGen.GenerateID()
	_, err := r.collection.InsertOne(ctx, player)
	if err != nil {
		player.ID = ""
	}
	return err
}

func (r *mongoPlayerRepository) FindByID(id string) (domain.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var player domain.Player
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	if mongo.ErrNoDocuments == err {
		return player, nil
	}
	return player, err
}

func (r *mongoPlayerRepository) FindByEmail(email string) (domain.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var player domain.Player
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&player)
	if mongo.ErrNoDocuments == err {
		return player, nil
	}
	return player, err
}

func (r *mongoPlayerRepository) FindByLastName(lastName string) ([]domain.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Using only json must be "all" lower case as mongo stores it in lower case, even if in the vew is showed in camel case;
	// but using bson then camel case is possible to use so enabling struct to support both is the right option,
	// i.e. using json fopr ReST and bson for mongo.
	//
	// Conclusion:
	// Using bson tags for MongoDB and json tags for HTTP APIs is a common and effective practice in Go applications.
	// This approach allows you to take advantage of MongoDB’s capabilities while seamlessly interacting with HTTP clients using JSON.
	// Just ensure you handle the conversion as needed between the two formats.
	cursor, err := r.collection.Find(ctx, bson.M{"lastName": lastName})
	if err != nil && mongo.ErrNoDocuments != err {
		return nil, err
	}
	defer cursor.Close(ctx)
	var players []domain.Player
	for cursor.Next(ctx) {
		var player domain.Player
		if err := cursor.Decode(&player); err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

func (r *mongoPlayerRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func NewMongoPlayerCoupleRepository(idGen utils.IDGenerator, client common.MongoClient) domain.PlayerCoupleRepository {
	collection := client.GetCollection(playerCouplesColName)
	return &mongoPlayerCoupleRepository{
		idGen:      idGen,
		collection: collection,
	}
}

func (r *mongoPlayerCoupleRepository) Upsert(playerCouple *domain.PlayerCouple) error {
	if playerCouple == nil {
		return errors.New("playerCouple is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// DDD repository principle.
	if len(playerCouple.ID) > 0 {
		_, err := r.collection.UpdateOne(ctx, bson.M{"_id": playerCouple.ID}, bson.M{"$set": playerCouple})
		return err
	}
	playerCouple.ID = r.idGen.GenerateIDWithPrefixes(playerCouple.Player1.LastName, playerCouple.Player2.LastName)
	_, err := r.collection.InsertOne(ctx, playerCouple)
	if err != nil {
		playerCouple.ID = ""
	}
	return err
}

func (r *mongoPlayerCoupleRepository) FindByID(id string) (domain.PlayerCouple, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var playerCouple domain.PlayerCouple
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&playerCouple)
	if mongo.ErrNoDocuments == err {
		return playerCouple, nil
	}
	return playerCouple, err
}

func (r *mongoPlayerCoupleRepository) FindByPrefixes(lastNamePlayer1, lastNamePlayer2 string) ([]domain.PlayerCouple, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var prefix = fmt.Sprintf("%s-%s", lastNamePlayer1, lastNamePlayer2)
	cursor, err := r.collection.Find(ctx, bson.M{
		"_id": bson.M{"$regex": "^" + prefix},
	})
	if err != nil && mongo.ErrNoDocuments != err {
		return nil, err
	}
	defer cursor.Close(ctx)

	var playerCouples []domain.PlayerCouple
	for cursor.Next(ctx) {
		var playerCouple domain.PlayerCouple
		if err := cursor.Decode(&playerCouple); err != nil {
			return nil, err
		}
		playerCouples = append(playerCouples, playerCouple)
	}
	return playerCouples, nil
}

func (r *mongoPlayerCoupleRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
