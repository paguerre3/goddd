package application

import "github.com/paguerre3/goddd/modules/player-couple/domain"

type playerCoupleService struct {
	playerRepo       domain.PlayerRepository
	playerCoupleRepo domain.PlayerCoupleRepository
	idGen            domain.IDGenerator
}
