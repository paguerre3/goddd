package application

import "github.com/paguerre3/goddd/modules/player-couple/domain"

type UnregisterPlayerUseCase interface {
	UnregisterPlayerUseCase(playerId string) (status UnregisterPlayerStatus, err error)
}

type UnregisterPlayerStatus uint8

const (
	UnregisterPlayerPending UnregisterPlayerStatus = iota
	UnregisterPlayerNotFound
	UnregisterPlayerDeleted
)

func NewUnregisterPlayerUseCase(playerRepository domain.PlayerRepository,
	idGenerator domain.IDGenerator) UnregisterPlayerUseCase {
	return &playerService{
		playerRepo: playerRepository,
		idGen:      idGenerator,
	}
}

func (s *playerService) UnregisterPlayerUseCase(playerId string) (status UnregisterPlayerStatus, err error) {
	foundPlayer, err := s.playerRepo.FindByID(playerId)
	if err != nil {
		return status, err
	}
	if len(foundPlayer.ID) == 0 {
		status = UnregisterPlayerNotFound
		return status, nil
	}
	if err = s.playerRepo.Delete(playerId); err != nil {
		return status, err
	}
	status = UnregisterPlayerDeleted
	return status, nil
}
