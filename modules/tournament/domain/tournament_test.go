package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock IDGenerator for testing
type MockIDGenerator struct{}

func (m *MockIDGenerator) GenerateID() string {
	return "mock-id"
}

func (m *MockIDGenerator) GenerateIDWithPrefixes(p1, p2 string) string {
	// empty implementtaion as it isn't used inside tournament:
	return ""
}

func TestTournament_MarshalJSON_Success(t *testing.T) {
	tournament := Tournament{
		ID:    "t1",
		Title: "Grand Slam",
		// "New" function isn't used so timestamp isn't validated over the time:
		Timestamp: time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC),
	}

	jsonData, err := json.Marshal(tournament)
	assert.NoError(t, err, "Expected no error during JSON marshalling")
	// Empty Players and Rounds slices should be omitted:
	assert.JSONEq(t, `{"id":"t1","title":"Grand Slam","timestamp":"2024-09-18T12:00"}`, string(jsonData), "Expected JSON to match")
}

func TestMatch_MarshalJSON_Success(t *testing.T) {
	match := Match{
		ID: "m1",
		// "New" function isn't used so timestamp isn't validated over the time:
		Timestamp: time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC),
	}

	jsonData, err := json.Marshal(match)
	assert.NoError(t, err, "Expected no error during JSON marshalling")
	// couple 1 and 2 are required so they can't be nil, instead they will always be empty if not created properly:
	assert.JSONEq(t, `{"id":"m1","timestamp":"2024-09-18T12:00","couple1":{"id":"","player1":{"id":"","email":"","name":"","surname":""},"player2":{"id":"","email":"","name":"","surname":""}},"couple2":{"id":"","player1":{"id":"","email":"","name":"","surname":""},"player2":{"id":"","email":"","name":"","surname":""}}}`,
		string(jsonData), "Expected JSON to match")
}

func TestNewTournament_Success(t *testing.T) {
	idGen := &MockIDGenerator{}
	timestamp := time.Now()
	playerCouples := []PlayerCouple{}
	rounds := []Round{}

	tournament, err := NewTournament(idGen, "Grand Slam", timestamp, playerCouples, rounds)
	assert.NoError(t, err, "Expected no error when creating a valid Tournament")
	assert.NotNil(t, tournament, "Expected Tournament to be non-nil")
}

func TestNewTournament_Fail_NilIDGenerator(t *testing.T) {
	timestamp := time.Now()
	playerCouples := []PlayerCouple{}
	rounds := []Round{}

	tournament, err := NewTournament(nil, "Grand Slam", timestamp, playerCouples, rounds)
	assert.Error(t, err, "Expected error when idGen is nil")
	assert.Nil(t, tournament, "Expected Tournament to be nil when idGen is nil")
}

func TestNewTournament_Fail_InvalidTitle(t *testing.T) {
	idGen := &MockIDGenerator{}
	timestamp := time.Now()
	playerCouples := []PlayerCouple{}
	rounds := []Round{}

	tournament, err := NewTournament(idGen, "GS", timestamp, playerCouples, rounds)
	assert.Error(t, err, "Expected error when title is too short")
	assert.Nil(t, tournament, "Expected Tournament to be nil when title is too short")
}

func TestNewTournament_Fail_OldTimestamp(t *testing.T) {
	idGen := &MockIDGenerator{}
	timestamp := time.Now().AddDate(0, 0, minMatchDays-1)
	playerCouples := []PlayerCouple{}
	rounds := []Round{}

	tournament, err := NewTournament(idGen, "Grand Slam", timestamp, playerCouples, rounds)
	assert.Error(t, err, "Expected error when timestamp is older than allowed")
	assert.Nil(t, tournament, "Expected Tournament to be nil when timestamp is older than allowed")
}

func TestNewRound_Success(t *testing.T) {
	matches := []Match{}

	round, err := NewRound(1, matches)
	assert.NoError(t, err, "Expected no error when creating a valid Round")
	assert.NotNil(t, round, "Expected Round to be non-nil")
}

func TestNewRound_Fail_InvalidRoundNumber(t *testing.T) {
	matches := []Match{}

	round, err := NewRound(-1, matches)
	assert.Error(t, err, "Expected error when roundNumber is invalid")
	assert.Nil(t, round, "Expected Round to be nil when roundNumber is invalid")
}

func TestNewMatch_Success(t *testing.T) {
	idGen := &MockIDGenerator{}
	timestamp := time.Now()
	couple1 := PlayerCouple{ID: "c1"}
	couple2 := PlayerCouple{ID: "c2"}
	score := &Score{}

	match, err := NewMatch(idGen, timestamp, couple1, couple2, score)
	assert.NoError(t, err, "Expected no error when creating a valid Match")
	assert.NotNil(t, match, "Expected Match to be non-nil")
}

func TestNewMatch_Fail_NilIDGenerator(t *testing.T) {
	timestamp := time.Now()
	couple1 := PlayerCouple{ID: "c1"}
	couple2 := PlayerCouple{ID: "c2"}
	score := &Score{}

	match, err := NewMatch(nil, timestamp, couple1, couple2, score)
	assert.Error(t, err, "Expected error when idGen is nil")
	assert.Nil(t, match, "Expected Match to be nil when idGen is nil")
}

func TestNewMatch_Fail_OldTimestamp(t *testing.T) {
	idGen := &MockIDGenerator{}
	timestamp := time.Now().AddDate(0, 0, minMatchDays-1)
	couple1 := PlayerCouple{ID: "c1"}
	couple2 := PlayerCouple{ID: "c2"}
	score := &Score{}

	match, err := NewMatch(idGen, timestamp, couple1, couple2, score)
	assert.Error(t, err, "Expected error when timestamp is older than allowed")
	assert.Nil(t, match, "Expected Match to be nil when timestamp is older than allowed")
}

func TestNewMatch_Fail_SameCouples(t *testing.T) {
	idGen := &MockIDGenerator{}
	timestamp := time.Now()
	couple1 := PlayerCouple{ID: "c1"}
	couple2 := PlayerCouple{ID: "c1"}
	score := &Score{}

	match, err := NewMatch(idGen, timestamp, couple1, couple2, score)
	assert.Error(t, err, "Expected error when couple1 and couple2 are the same")
	assert.Nil(t, match, "Expected Match to be nil when couple1 and couple2 are the same")
}

func TestNewScore_Success(t *testing.T) {
	set1 := &GameSet{GamesCouple1: 6, GamesCouple2: 4}
	set2 := &GameSet{GamesCouple1: 3, GamesCouple2: 6}
	set3 := &GameSet{GamesCouple1: 6, GamesCouple2: 2}

	score, err := NewScore(set1, set2, set3)
	assert.NoError(t, err, "Expected no error when creating a valid Score")
	assert.NotNil(t, score, "Expected Score to be non-nil")
}

func TestNewScore_Fail_NilSet1(t *testing.T) {
	set2 := &GameSet{GamesCouple1: 6, GamesCouple2: 3}

	score, err := NewScore(nil, set2, nil)
	assert.Error(t, err, "Expected error when set1 is nil")
	assert.Nil(t, score, "Expected Score to be nil when set1 is nil")
}

func TestNewScore_Fail_NilSet2(t *testing.T) {
	set1 := &GameSet{GamesCouple1: 6, GamesCouple2: 4}

	score, err := NewScore(set1, nil, nil)
	assert.Error(t, err, "Expected error when set2 is nil")
	assert.Nil(t, score, "Expected Score to be nil when set2 is nil")
}

func TestNewGameSet_Success(t *testing.T) {
	tiebreak := &Tiebreak{PointsCouple1: 7, PointsCouple2: 5}
	gameSet, err := NewGameSet(6, 6, tiebreak)
	assert.NoError(t, err, "Expected no error when creating a valid GameSet")
	assert.NotNil(t, gameSet, "Expected GameSet to be non-nil")
}

func TestNewGameSet_Fail_InvalidGamesCouple1(t *testing.T) {
	gameSet, err := NewGameSet(-1, 4, nil)
	assert.Error(t, err, "Expected error when gamesCouple1 is invalid")
	assert.Nil(t, gameSet, "Expected GameSet to be nil when gamesCouple1 is invalid")
}

func TestNewGameSet_Fail_InvalidGamesCouple2(t *testing.T) {
	tiebreak := &Tiebreak{PointsCouple1: 7, PointsCouple2: 5}
	gameSet, err := NewGameSet(6, 20, tiebreak)
	assert.Error(t, err, "Expected error when gamesCouple2 is invalid")
	assert.Nil(t, gameSet, "Expected GameSet to be nil when gamesCouple2 is invalid")
}

func TestNewTiebreak_Success(t *testing.T) {
	tiebreak, err := NewTiebreak(7, 5)
	assert.NoError(t, err, "Expected no error when creating a valid Tiebreak")
	assert.NotNil(t, tiebreak, "Expected Tiebreak to be non-nil")
}

func TestNewTiebreak_Fail_InvalidPointsCouple1(t *testing.T) {
	tiebreak, err := NewTiebreak(-1, 1)
	assert.Error(t, err, "Expected error when pointsCouple1 is invalid")
	assert.Nil(t, tiebreak, "Expected Tiebreak to be nil when pointsCouple1 is invalid")
}

func TestNewTiebreak_Fail_InvalidPointsCouple2(t *testing.T) {
	tiebreak, err := NewTiebreak(19, 22)
	assert.Error(t, err, "Expected error when pointsCouple2 is invalid")
	assert.Nil(t, tiebreak, "Expected Tiebreak to be nil when pointsCouple2 is invalid")
}
