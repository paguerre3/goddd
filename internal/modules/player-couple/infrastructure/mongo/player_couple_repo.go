package mongo

import (
	"context"
	"fmt"
	"time"

	common "github.com/paguerre3/goddd/internal/modules/common/mongo"
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
	collection *mongo.Collection
}

type mongoPlayerCoupleRepository struct {
	collection *mongo.Collection
}

func NewMongoPlayerRepository(client common.MongoClient) domain.PlayerRepository {
	collection := client.GetCollection(playersColName)
	return &mongoPlayerRepository{collection: collection}
}

func (r *mongoPlayerRepository) Save(player domain.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, player)
	return err
}

func (r *mongoPlayerRepository) FindByID(id string) (domain.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var player domain.Player
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&player)
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

	// must be "all" lower case as mongo stores it in lower case, even if in the vew is showed in camel case:
	cursor, err := r.collection.Find(ctx, bson.M{"lastname": lastName})
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

func (r *mongoPlayerRepository) Update(player domain.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"id": player.ID}, bson.M{"$set": player})
	return err
}

func (r *mongoPlayerRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func NewMongoPlayerCoupleRepository(client common.MongoClient) domain.PlayerCoupleRepository {
	collection := client.GetCollection(playerCouplesColName)
	return &mongoPlayerCoupleRepository{collection: collection}
}

func (r *mongoPlayerCoupleRepository) Save(playerCouple domain.PlayerCouple) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, playerCouple)
	return err
}

func (r *mongoPlayerCoupleRepository) FindByID(id string) (domain.PlayerCouple, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var playerCouple domain.PlayerCouple
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&playerCouple)
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
		"id": bson.M{"$regex": "^" + prefix},
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

func (r *mongoPlayerCoupleRepository) Update(playerCouple domain.PlayerCouple) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"id": playerCouple.ID}, bson.M{"$set": playerCouple})
	return err
}

func (r *mongoPlayerCoupleRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
