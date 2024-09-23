package domain

import (
	"fmt"
	"net/mail"
)

const (
	minRank       = 1
	maxRank       = 8
	minSSNDigits  = 8
	minAge        = 3
	maxAge        = 100
	minNameDigits = 3
)

type Player struct {
	ID                   string  `json:"id"`
	Email                string  `json:"email"`
	SocialSecurityNumber *string `json:"ssn,omitempty"`
	FirstName            string  `json:"name"`
	LastName             string  `json:"surname"`
	Age                  *int    `json:"age,omitempty"`
}

type PlayerCouple struct {
	ID      string `json:"id"`
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Ranking *int   `json:"ranking,omitempty"`
}

func NewPlayer(idGen IDGenerator, email string, socialSecurityNumber *string,
	firstName, lastName string, age *int) (*Player, error) {
	if idGen == nil {
		return nil, fmt.Errorf("idGen cannot be nil")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("invalid email: %s", email)
	}
	if socialSecurityNumber != nil && len(*socialSecurityNumber) < minSSNDigits {
		return nil, fmt.Errorf("invalid social security number: %s", *socialSecurityNumber)
	}
	if len(firstName) < minNameDigits {
		return nil, fmt.Errorf("invalid first name: %s", firstName)
	}
	if len(lastName) < minNameDigits {
		return nil, fmt.Errorf("invalid last name: %s", lastName)
	}
	if age != nil && (*age < minAge || *age > maxAge) {
		return nil, fmt.Errorf("invalid age: %d", *age)
	}
	return &Player{
		ID:                   idGen.GenerateID(),
		Email:                email,
		SocialSecurityNumber: socialSecurityNumber,
		FirstName:            firstName,
		LastName:             lastName,
		Age:                  age,
	}, nil
}

func NewPlayerCouple(idGen IDGenerator, player1 Player, player2 Player, ranking *int) (*PlayerCouple, error) {
	if idGen == nil {
		return nil, fmt.Errorf("idGen cannot be nil")
	}
	if player1.ID == player2.ID {
		return nil, fmt.Errorf("player1 and player2 cannot be the same")
	}
	if ranking != nil && (*ranking < minRank || *ranking > maxRank) {
		return nil, fmt.Errorf("invalid ranking: %d", *ranking)
	}
	return &PlayerCouple{
		ID:      idGen.GenerateIDWithPrefixes(player1.LastName, player2.LastName),
		Player1: player1,
		Player2: player2,
		Ranking: ranking,
	}, nil
}
