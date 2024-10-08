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
	minIdDigits   = 3
)

type Player struct {
	ID                   string  `bson:"_id" json:"id"`
	Email                string  `bson:"email" json:"email"`
	SocialSecurityNumber *string `bson:"socialSecurityNumber,omitempty" json:"socialSecurityNumber,omitempty"`
	FirstName            string  `bson:"firstName" json:"firstName"`
	LastName             string  `bson:"lastName" json:"lastName"`
	Age                  *int    `bson:"age,omitempty" json:"age,omitempty"`
}

type PlayerCouple struct {
	ID      string `bson:"_id" json:"id"`
	Player1 Player `bson:"player1" json:"player1"`
	Player2 Player `bson:"player2" json:"player2"`
	Ranking *int   `bson:"ranking,omitempty" json:"ranking,omitempty"`
}

func NewPlayer(email string, socialSecurityNumber *string,
	firstName, lastName string, age *int) (*Player, error) {
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if socialSecurityNumber != nil && len(*socialSecurityNumber) < minSSNDigits {
		return nil, fmt.Errorf("invalid social security number: %s", *socialSecurityNumber)
	}
	if len(firstName) < minNameDigits {
		return nil, fmt.Errorf("invalid first name: %s", firstName)
	}
	if err := ValidateLastName(lastName); err != nil {
		return nil, err
	}
	if age != nil && (*age < minAge || *age > maxAge) {
		return nil, fmt.Errorf("invalid age: %d", *age)
	}
	return &Player{
		//ID:                 auto generated ID set in the repository.
		Email:                email,
		SocialSecurityNumber: socialSecurityNumber,
		FirstName:            firstName,
		LastName:             lastName,
		Age:                  age,
	}, nil
}

func ValidateID(id string) error {
	if len(id) < minIdDigits {
		return fmt.Errorf("invalid id: %s", id)
	}
	return nil
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email: %s", email)
	}
	return nil
}

func ValidateLastName(lastName string) error {
	if len(lastName) < minNameDigits {
		return fmt.Errorf("invalid last name: %s", lastName)
	}
	return nil
}

func NewPlayerCouple(player1 Player, player2 Player, ranking *int) (*PlayerCouple, error) {
	if player1.ID == player2.ID {
		return nil, fmt.Errorf("player1 and player2 cannot be the same")
	}
	if ranking != nil && (*ranking < minRank || *ranking > maxRank) {
		return nil, fmt.Errorf("invalid ranking: %d", *ranking)
	}
	return &PlayerCouple{
		//ID:    auto generated ID set in the repository.
		Player1: player1,
		Player2: player2,
		Ranking: ranking,
	}, nil
}
