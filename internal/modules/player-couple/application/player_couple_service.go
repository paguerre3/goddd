package application

import "github.com/paguerre3/goddd/internal/modules/player-couple/domain"

type playerService struct {
	playerRepo domain.PlayerRepository
	idGen      domain.IDGenerator
}
