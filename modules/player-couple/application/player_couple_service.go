package application

import "github.com/paguerre3/goddd/modules/player-couple/domain"

type PlayerCoupleService struct {
	playerRepo       domain.PlayerRepository
	playerCoupleRepo domain.PlayerCoupleRepository
	idGen            domain.IDGenerator
}

func NewPlayerCoupleService(playerRepository domain.PlayerRepository, playerCoupleRepository domain.PlayerCoupleRepository,
	idGenerator domain.IDGenerator) *PlayerCoupleService {
	return &PlayerCoupleService{
		playerRepo:       playerRepository,
		playerCoupleRepo: playerCoupleRepository,
		idGen:            idGenerator,
	}
}
