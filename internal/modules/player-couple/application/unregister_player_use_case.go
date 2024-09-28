package application

import "github.com/paguerre3/goddd/internal/modules/player-couple/domain"

type UnregisterPlayerUseCase interface {
	UnregisterPlayerUseCase(playerId string) (status UnregisterPlayerStatus, err error)
}

type UnregisterPlayerStatus uint8

const (
	UnregisterPlayerPending UnregisterPlayerStatus = iota
	UnregisterPlayerInvalid
	UnregisterPlayerNotFound
	UnregisterPlayerDeleted
)

// Implement the Stringer interface.
func (s UnregisterPlayerStatus) String() string {
	return [...]string{"UnregisterPlayerPending", "UnregisterPlayerInvalid", "UnregisterPlayerNotFound", "UnregisterPlayerDeleted"}[s]
}

func NewUnregisterPlayerUseCase(idGenerator domain.IDGenerator,
	playerRepository domain.PlayerRepository) UnregisterPlayerUseCase {
	return &playerService{
		idGen:      idGenerator,
		playerRepo: playerRepository,
	}
}

func (s *playerService) UnregisterPlayerUseCase(playerId string) (status UnregisterPlayerStatus, err error) {
	if err := domain.ValidateID(playerId); err != nil {
		status = UnregisterPlayerInvalid
		return status, err
	}
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
