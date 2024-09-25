package application

import "github.com/paguerre3/goddd/modules/player-couple/domain"

type PlayerService struct {
	playerRepo domain.PlayerRepository
	idGen      domain.IDGenerator
}

type PlayerCoupleService struct {
	playerRepo       domain.PlayerRepository
	playerCoupleRepo domain.PlayerCoupleRepository
	idGen            domain.IDGenerator
}

func NewPlayerService(playerRepository domain.PlayerRepository, idGenerator domain.IDGenerator) *PlayerService {
	return &PlayerService{playerRepo: playerRepository, idGen: idGenerator}
}

func NewPlayerCoupleService(playerRepository domain.PlayerRepository, playerCoupleRepository domain.PlayerCoupleRepository,
	idGenerator domain.IDGenerator) *PlayerCoupleService {
	return &PlayerCoupleService{
		playerRepo:       playerRepository,
		playerCoupleRepo: playerCoupleRepository,
		idGen:            idGenerator,
	}
}
