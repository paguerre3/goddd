package domain

// interfaces to be used by infrastructure layer:
type PlayerRepository interface {
	Save(player Player) error
	FindByID(id string) (Player, error)
	FindByEmail(email string) (Player, error)
	FindByLastName(lastName string) ([]Player, error)
	Update(player Player) error
	Delete(id string) error
}

type PlayerCoupleRepository interface {
	Save(playerCouple PlayerCouple) error
	FindByID(id string) (PlayerCouple, error)
	FindByPrefixes(lastNamePlayer1, lastNamePlayer2 string) ([]PlayerCouple, error)
	Update(playerCouple PlayerCouple) error
	Delete(id string) error
}
