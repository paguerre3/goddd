package domain

// Adapted copy similar to player_couple module.
// (not imported as its a different module so it must be a copy for tournaments).
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
