package domain

// Adapted copy similar to player_couple module.
// (not imported as its a different module so it must be a copy for tournaments).
// Used to retrieve registered players domain objects from tournament (repository).
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
