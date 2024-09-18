package domain

import "fmt"

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
	Ranking int    `json:"ranking"`
}

func NewPlayer(socialSecurityNumber *string, name string, age *int) (*Player, error) {
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
		ID:   "", // TODO: generate ID
		Name: name,
		Age:  age,
	}, nil
}

func NewPlayerCouple(id string, player1 Player, player2 Player, ranking int) *PlayerCouple {
	return &PlayerCouple{
		ID:      id,
		Player1: player1,
		Player2: player2,
		Ranking: ranking,
	}
}
