package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockId = "mock-id"
)

// Mock IDGenerator for testing without breaking modularity principle even in testing domain module package.
type MockIDGenerator struct {
	// molck aggregate used as a helper mock to avoid repetitions:
	aggregate *string
}

func (m *MockIDGenerator) GenerateID() string {
	if m.aggregate != nil {
		return fmt.Sprintf("%s-%s", mockId, *m.aggregate)
	}
	return mockId
}

func (m *MockIDGenerator) GenerateIDWithPrefixes(p1, p2 string) string {
	return fmt.Sprintf("%s&%s-%s", p1, p2, mockId)
}

// TestNewPlayer_Success tests successful creation of a Player
func TestNewPlayer_Success(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 25
	player, err := NewPlayer(idGen, &ssn, "Roger Federer", &age)

	assert.NoError(t, err, "Expected no error for valid inputs")
	assert.NotNil(t, player, "Expected player object to be created")
	assert.Equal(t, "mock-id", player.ID, "Expected ID to be mock-id")
	assert.Equal(t, &ssn, player.SocialSecurityNumber, "Expected correct SSN")
	assert.Equal(t, "Roger Federer", player.Name, "Expected correct player name")
	assert.Equal(t, &age, player.Age, "Expected correct age")
}

// TestNewPlayer_Fail_InvalidSSN tests failure for invalid SSN
func TestNewPlayer_Fail_InvalidSSN(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "123"
	age := 25
	player, err := NewPlayer(idGen, &ssn, "Roger Federer", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid SSN")
	assert.EqualError(t, err, "invalid social security number: 123", "Expected invalid SSN error")
}

// TestNewPlayer_Fail_InvalidName tests failure for invalid player name
func TestNewPlayer_Fail_InvalidName(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 25
	player, err := NewPlayer(idGen, &ssn, "R", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid name")
	assert.EqualError(t, err, "invalid name: R", "Expected invalid name error")
}

// TestNewPlayer_Fail_InvalidAge tests failure for invalid age
func TestNewPlayer_Fail_InvalidAge(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 120
	player, err := NewPlayer(idGen, &ssn, "Roger Federer", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid age")
	assert.EqualError(t, err, "invalid age: 120", "Expected invalid age error")
}

// TestNewPlayerCouple_Success tests successful creation of a PlayerCouple
func TestNewPlayerCouple_Success(t *testing.T) {
	agg1 := "a"
	agg2 := "b"
	idGen := &MockIDGenerator{aggregate: &agg1}
	// Creating two valid players
	player1, _ := NewPlayer(idGen, nil, "Roger Federer", nil)
	idGen.aggregate = &agg2
	player2, _ := NewPlayer(idGen, nil, "Rafael Nadal", nil)

	couple, err := NewPlayerCouple(idGen, *player1, *player2, nil)

	assert.NoError(t, err, "Expected no error for valid couple creation")
	assert.NotNil(t, couple, "Expected couple object to be created")
	assert.Equal(t, "Roger Federer&Rafael Nadal-mock-id", couple.ID, "Expected correct couple ID")
}

// TestNewPlayerCouple_Fail_SamePlayer tests failure when Player1 and Player2 are the same
func TestNewPlayerCouple_Fail_SamePlayer(t *testing.T) {
	idGen := &MockIDGenerator{}

	player, _ := NewPlayer(idGen, nil, "Roger Federer", nil)

	couple, err := NewPlayerCouple(idGen, *player, *player, nil)

	assert.Nil(t, couple, "Expected no couple to be created when Player1 and Player2 are the same")
	assert.EqualError(t, err, "player1 and player2 cannot be the same", "Expected player1 and player2 cannot be the same error")
}

// TestNewPlayerCouple_Fail_InvalidRanking tests failure for invalid ranking
func TestNewPlayerCouple_Fail_InvalidRanking(t *testing.T) {
	agg1 := "a"
	agg2 := "b"
	idGen := &MockIDGenerator{aggregate: &agg1}

	player1, _ := NewPlayer(idGen, nil, "Roger Federer", nil)
	idGen.aggregate = &agg2
	player2, _ := NewPlayer(idGen, nil, "Rafael Nadal", nil)
	ranking := 9

	couple, err := NewPlayerCouple(idGen, *player1, *player2, &ranking)

	assert.Nil(t, couple, "Expected no couple to be created with invalid ranking")
	assert.EqualError(t, err, "invalid ranking: 9", "Expected invalid ranking error")
}
