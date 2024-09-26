package domain

// Adapted copy similar to player_couple module.
// (not imported as its a different module so it must be a copy for tournaments).
// Used to retrieve registered players domain objects from tournament (repository).
type Player struct {
	ID                   string  `json:"id"`
	Email                string  `json:"email"`
	SocialSecurityNumber *string `json:"socialSecurityNumber,omitempty"`
	FirstName            string  `json:"firstName"`
	LastName             string  `json:"lastName"`
	Age                  *int    `json:"age,omitempty"`
}

type PlayerCouple struct {
	ID      string `json:"id"`
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Ranking *int   `json:"ranking,omitempty"`
}
