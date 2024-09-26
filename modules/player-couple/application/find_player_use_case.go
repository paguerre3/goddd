package application

import "github.com/paguerre3/goddd/modules/player-couple/domain"

type FindPlayerUseCase interface {
	FindPlayerByIDUseCase(playerId string) (domain.Player, FindPlayerStatus, error)
	FindPlayerByEmailUseCase(email string) (domain.Player, FindPlayerStatus, error)
	FindPlayersByLastNameUseCase(lastName string) ([]domain.Player, FindPlayerStatus, error)
}

type FindPlayerStatus uint8

const (
	FindPlayerPending FindPlayerStatus = iota
	FindPlayerInvalid
	FindPlayerNotFound
	FindPlayerFound
)

func NewFindPlayerUseCase(playerRepository domain.PlayerRepository,
	idGenerator domain.IDGenerator) FindPlayerUseCase {
	return &playerService{
		playerRepo: playerRepository,
		idGen:      idGenerator,
	}
}

func (s *playerService) FindPlayerByIDUseCase(playerId string) (domain.Player, FindPlayerStatus, error) {
	if err := domain.ValidateID(playerId); err != nil {
		return domain.Player{}, FindPlayerInvalid, err
	}
	player, err := s.playerRepo.FindByID(playerId)
	if err != nil {
		return player, FindPlayerPending, err
	}
	if len(player.ID) == 0 {
		return player, FindPlayerNotFound, nil
	}
	return player, FindPlayerFound, nil
}

func (s *playerService) FindPlayerByEmailUseCase(email string) (domain.Player, FindPlayerStatus, error) {
	if err := domain.ValidateEmail(email); err != nil {
		return domain.Player{}, FindPlayerInvalid, err
	}
	player, err := s.playerRepo.FindByEmail(email)
	if err != nil {
		return player, FindPlayerPending, err
	}
	if len(player.ID) == 0 {
		return player, FindPlayerNotFound, nil
	}
	return player, FindPlayerFound, nil
}

func (s *playerService) FindPlayersByLastNameUseCase(lastName string) ([]domain.Player, FindPlayerStatus, error) {
	if err := domain.ValidateLastName(lastName); err != nil {
		return nil, FindPlayerInvalid, err
	}
	players, err := s.playerRepo.FindByLastName(lastName)
	if err != nil {
		return players, FindPlayerPending, err
	}
	if len(players) == 0 {
		return players, FindPlayerNotFound, nil
	}
	return players, FindPlayerFound, nil
}
