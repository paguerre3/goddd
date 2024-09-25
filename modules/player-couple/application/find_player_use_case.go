package application

import "github.com/paguerre3/goddd/modules/player-couple/domain"

type FindPlayerUseCase interface {
	FindPlayerByIDUseCase(playerId string) (domain.Player, error)
	FindPlayerByEmailUseCase(email string) (domain.Player, error)
	FindPlayersByLastNameUseCase(lastName string) ([]domain.Player, error)
}

func NewFindPlayerUseCase(playerRepository domain.PlayerRepository,
	idGenerator domain.IDGenerator) FindPlayerUseCase {
	return &playerService{
		playerRepo: playerRepository,
		idGen:      idGenerator,
	}
}

func (s *playerService) FindPlayerByIDUseCase(playerId string) (domain.Player, error) {
	if err := domain.ValidateID(playerId); err != nil {
		return domain.Player{}, err
	}
	return s.playerRepo.FindByID(playerId)
}

func (s *playerService) FindPlayerByEmailUseCase(email string) (domain.Player, error) {
	if err := domain.ValidateEmail(email); err != nil {
		return domain.Player{}, err
	}
	return s.playerRepo.FindByEmail(email)
}

func (s *playerService) FindPlayersByLastNameUseCase(lastName string) ([]domain.Player, error) {
	if err := domain.ValidateLastName(lastName); err != nil {
		return nil, err
	}
	return s.playerRepo.FindByLastName(lastName)
}
