package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockId        = "mock-id"
	anotherMockId = "another-mock-id"
)

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
	ssn := "12345678"
	age := 25
	player, err := NewPlayer("agus.tapia@gmail.com", &ssn, "Agustin", "Tapia", &age)

	assert.NoError(t, err, "Expected no error for valid inputs")
	assert.NotNil(t, player, "Expected player object to be created")
	assert.Equal(t, "", player.ID, "Expected ID to be empty")
	assert.Equal(t, "agus.tapia@gmail.com", player.Email, "Expected correct email")
	assert.Equal(t, &ssn, player.SocialSecurityNumber, "Expected correct SSN")
	assert.Equal(t, "Agustin", player.FirstName, "Expected correct player name")
	assert.Equal(t, "Tapia", player.LastName, "Expected correct player surname")
	assert.Equal(t, &age, player.Age, "Expected correct age")
}

// TestNewPlayer_Fail_InvalidEmail tests failure for invalid email
func TestNewPlayer_Fail_InvalidEmail(t *testing.T) {
	ssn := "12345678"
	age := 25
	player, err := NewPlayer("agustapia", &ssn, "Agustin", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid email")
	assert.EqualError(t, err, "invalid email: agustapia", "Expected invalid email")
}

// TestNewPlayer_Fail_InvalidSSN tests failure for invalid SSN
func TestNewPlayer_Fail_InvalidSSN(t *testing.T) {
	ssn := "123"
	age := 25
	player, err := NewPlayer("agus.tapia@gmail.com", &ssn, "Agustin", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid SSN")
	assert.EqualError(t, err, "invalid social security number: 123", "Expected invalid SSN error")
}

// TestNewPlayer_Fail_InvalidFirstName tests failure for invalid player first name
func TestNewPlayer_Fail_InvalidFirstName(t *testing.T) {
	ssn := "12345678"
	age := 25
	player, err := NewPlayer("agus.tapia@gmail.com", &ssn, "A", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid first name")
	assert.EqualError(t, err, "invalid first name: A", "Expected invalid first name error")
}

// TestNewPlayer_Fail_InvalidLastName tests failure for invalid player last name
func TestNewPlayer_Fail_InvalidLastName(t *testing.T) {
	ssn := "12345678"
	age := 25
	player, err := NewPlayer("agus.tapia@gmail.com", &ssn, "Agustin", "T", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid last name")
	assert.EqualError(t, err, "invalid last name: T", "Expected invalid last name error")
}

// TestNewPlayer_Fail_InvalidAge tests failure for invalid age
func TestNewPlayer_Fail_InvalidAge(t *testing.T) {
	ssn := "12345678"
	age := 120
	player, err := NewPlayer("agus.tapia@gmail.com", &ssn, "Agustin", "Tapia", &age)

	assert.Nil(t, player, "Expected no player to be created with invalid age")
	assert.EqualError(t, err, "invalid age: 120", "Expected invalid age error")
}

// TestNewPlayerCouple_Success tests successful creation of a PlayerCouple
func TestNewPlayerCouple_Success(t *testing.T) {
	// Creating two valid players
	player1, _ := NewPlayer("agus.tapia@gmail.com", nil, "Agustin", "Tapia", nil)
	player1.ID = mockId
	player2, _ := NewPlayer("ale.galan@gmail.com", nil, "Ale", "Galan", nil)
	player2.ID = anotherMockId

	couple, err := NewPlayerCouple(*player1, *player2, nil)
	// mock id generated in repository
	couple.ID = fmt.Sprintf("%s-%s-%s", player1.LastName, player2.LastName, mockId)

	assert.NoError(t, err, "Expected no error for valid couple creation")
	assert.NotNil(t, couple, "Expected couple object to be created")
	assert.Equal(t, "Tapia-Galan-mock-id", couple.ID, "Expected correct couple ID")
}

// TestNewPlayerCouple_Fail_SamePlayer tests failure when Player1 and Player2 are the same
func TestNewPlayerCouple_Fail_SamePlayer(t *testing.T) {
	player, _ := NewPlayer("agus.tapia@gmail.com", nil, "Agustin", "Tapia", nil)
	player.ID = mockId

	couple, err := NewPlayerCouple(*player, *player, nil)

	assert.Nil(t, couple, "Expected no couple to be created when Player1 and Player2 are the same")
	assert.EqualError(t, err, "player1 and player2 cannot be the same", "Expected player1 and player2 cannot be the same error")
}

// TestNewPlayerCouple_Fail_InvalidRanking tests failure for invalid ranking
func TestNewPlayerCouple_Fail_InvalidRanking(t *testing.T) {
	player1, _ := NewPlayer("agus.tapia@gmail.com", nil, "Agustin", "Tapia", nil)
	player1.ID = mockId
	player2, _ := NewPlayer("ale.galan@gmail.com", nil, "Ale", "Galan", nil)
	player2.ID = anotherMockId
	ranking := 9

	couple, err := NewPlayerCouple(*player1, *player2, &ranking)

	assert.Nil(t, couple, "Expected no couple to be created with invalid ranking")
	assert.EqualError(t, err, "invalid ranking: 9", "Expected invalid ranking error")
}
