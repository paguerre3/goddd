package domain

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type PlayerCouple struct {
	ID      string `json:"id"`
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Ranking int    `json:"ranking"`
}

func NewValidPlayer(id string, name string, age int) *Player {
	return &Player{
		ID:   id,
		Name: name,
		Age:  age,
	}
}

func NewValidPlayerCouple(id string, player1 Player, player2 Player, ranking int) *PlayerCouple {
	return &PlayerCouple{
		ID:      id,
		Player1: player1,
		Player2: player2,
		Ranking: ranking,
	}
}
