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
	return fmt.Sprintf("%s-%s-%s", p1, p2, mockId)
}

func TestValidateID(t *testing.T) {
	// Test case: ID is shorter than minIdDigits
	id := "12"
	err := ValidateID(id)
	assert.Error(t, err, "Expected error for ID shorter than minIdDigits")
	assert.EqualError(t, err, "invalid id: 12", "Expected error message for ID shorter than minIdDigits")

	// Test case: ID is equal to minIdDigits
	id = "123"
	err = ValidateID(id)
	assert.NoError(t, err, "Expected no error for ID equal to minIdDigits")

	// Test case: ID is longer than minIdDigits
	id = "123456789012345678901234567890"
	err = ValidateID(id)
	assert.NoError(t, err, "Expected no error for ID longer than minIdDigits")
}

// TestNewPlayer_Success tests successful creation of a Player
func TestNewPlayer_Success(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 25
	player, err := NewPlayer(idGen, "agus.tapia@gmail.com", &ssn, "Agustin", "Tapia", &age)

	assert.NoError(t, err, "Expected no error for valid inputs")
	assert.NotNil(t, player, "Expected player object to be created")
	assert.Equal(t, "mock-id", player.ID, "Expected ID to be mock-id")
	assert.Equal(t, "agus.tapia@gmail.com", player.Email, "Expected correct email")
	assert.Equal(t, &ssn, player.SocialSecurityNumber, "Expected correct SSN")
	assert.Equal(t, "Agustin", player.FirstName, "Expected correct player name")
	assert.Equal(t, "Tapia", player.LastName, "Expected correct player surname")
	assert.Equal(t, &age, player.Age, "Expected correct age")
}

// TestNewPlayer_Fail_InvalidEmail tests failure for invalid email
func TestNewPlayer_Fail_InvalidEmail(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 25
	player, err := NewPlayer(idGen, "agustapia", &ssn, "Agustin", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid email")
	assert.EqualError(t, err, "invalid email: agustapia", "Expected invalid email")
}

// TestNewPlayer_Fail_InvalidSSN tests failure for invalid SSN
func TestNewPlayer_Fail_InvalidSSN(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "123"
	age := 25
	player, err := NewPlayer(idGen, "agus.tapia@gmail.com", &ssn, "Agustin", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid SSN")
	assert.EqualError(t, err, "invalid social security number: 123", "Expected invalid SSN error")
}

// TestNewPlayer_Fail_InvalidFirstName tests failure for invalid player first name
func TestNewPlayer_Fail_InvalidFirstName(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 25
	player, err := NewPlayer(idGen, "agus.tapia@gmail.com", &ssn, "A", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid first name")
	assert.EqualError(t, err, "invalid first name: A", "Expected invalid first name error")
}

// TestNewPlayer_Fail_InvalidLastName tests failure for invalid player last name
func TestNewPlayer_Fail_InvalidLastName(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 25
	player, err := NewPlayer(idGen, "agus.tapia@gmail.com", &ssn, "Agustin", "T", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid last name")
	assert.EqualError(t, err, "invalid last name: T", "Expected invalid last name error")
}

// TestNewPlayer_Fail_InvalidAge tests failure for invalid age
func TestNewPlayer_Fail_InvalidAge(t *testing.T) {
	idGen := &MockIDGenerator{}
	ssn := "12345678"
	age := 120
	player, err := NewPlayer(idGen, "agus.tapia@gmail.com", &ssn, "Agustin", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid age")
	assert.EqualError(t, err, "invalid age: 120", "Expected invalid age error")
}

// TestNewPlayerCouple_Success tests successful creation of a PlayerCouple
func TestNewPlayerCouple_Success(t *testing.T) {
	agg1 := "a"
	agg2 := "b"
	idGen := &MockIDGenerator{aggregate: &agg1}
	// Creating two valid players
	player1, _ := NewPlayer(idGen, "agus.tapia@gmail.com", nil, "Agustin", "Tapia", nil)
	idGen.aggregate = &agg2
	player2, _ := NewPlayer(idGen, "ale.galan@gmail.com", nil, "Ale", "Galan", nil)

	couple, err := NewPlayerCouple(idGen, *player1, *player2, nil)

	assert.NoError(t, err, "Expected no error for valid couple creation")
	assert.NotNil(t, couple, "Expected couple object to be created")
	assert.Equal(t, "Tapia-Galan-mock-id", couple.ID, "Expected correct couple ID")
}

// TestNewPlayerCouple_Fail_SamePlayer tests failure when Player1 and Player2 are the same
func TestNewPlayerCouple_Fail_SamePlayer(t *testing.T) {
	idGen := &MockIDGenerator{}

	player, _ := NewPlayer(idGen, "agus.tapia@gmail.com", nil, "Agustin", "Tapia", nil)

	couple, err := NewPlayerCouple(idGen, *player, *player, nil)

	assert.Nil(t, couple, "Expected no couple to be created when Player1 and Player2 are the same")
	assert.EqualError(t, err, "player1 and player2 cannot be the same", "Expected player1 and player2 cannot be the same error")
}

// TestNewPlayerCouple_Fail_InvalidRanking tests failure for invalid ranking
func TestNewPlayerCouple_Fail_InvalidRanking(t *testing.T) {
	agg1 := "a"
	agg2 := "b"
	idGen := &MockIDGenerator{aggregate: &agg1}

	player1, _ := NewPlayer(idGen, "agus.tapia@gmail.com", nil, "Agustin", "Tapia", nil)
	idGen.aggregate = &agg2
	player2, _ := NewPlayer(idGen, "ale.galan@gmail.com", nil, "Ale", "Galan", nil)
	ranking := 9

	couple, err := NewPlayerCouple(idGen, *player1, *player2, &ranking)

	assert.Nil(t, couple, "Expected no couple to be created with invalid ranking")
	assert.EqualError(t, err, "invalid ranking: 9", "Expected invalid ranking error")
}
