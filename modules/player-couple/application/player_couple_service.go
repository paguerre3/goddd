package application

import "github.com/paguerre3/goddd/modules/player-couple/domain"

type playerService struct {
	playerRepo domain.PlayerRepository
	idGen      domain.IDGenerator
}
