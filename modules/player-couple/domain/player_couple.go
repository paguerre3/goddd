package domain

import (
	"fmt"
)

type Player struct {
	ID                   string  `json:"id"`
	SocialSecurityNumber *string `json:"ssn,omitempty"`
	Name                 string  `json:"name"`
	Age                  *int    `json:"age,omitempty"`
}

type PlayerCouple struct {
	ID      string `json:"id"`
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Ranking *int   `json:"ranking,omitempty"`
}

func NewPlayer(idGen IDGenerator, socialSecurityNumber *string,
	name string, age *int) (*Player, error) {
	if idGen == nil {
		return nil, fmt.Errorf("idGen cannot be nil")
	}
	if socialSecurityNumber != nil && len(*socialSecurityNumber) < 8 {
		return nil, fmt.Errorf("invalid social security number: %s", *socialSecurityNumber)
	}
	if len(name) < 2 {
		return nil, fmt.Errorf("invalid name: %s", name)
	}
	if age != nil && (*age < 3 || *age > 100) {
		return nil, fmt.Errorf("invalid age: %d", *age)
	}
	return &Player{
		ID:                   idGen.GenerateID(),
		SocialSecurityNumber: socialSecurityNumber,
		Name:                 name,
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
	if ranking != nil && (*ranking < 1 || *ranking > 8) {
		return nil, fmt.Errorf("invalid ranking: %d", *ranking)
	}
	return &PlayerCouple{
		ID:      idGen.GenerateIDWithPrefixes(player1.Name, player2.Name),
		Player1: player1,
		Player2: player2,
		Ranking: ranking,
	}, nil
}
