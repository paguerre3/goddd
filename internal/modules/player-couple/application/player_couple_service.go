package application

import "github.com/paguerre3/goddd/internal/modules/player-couple/domain"

type playerService struct {
	idGen      domain.IDGenerator
	playerRepo domain.PlayerRepository
}
