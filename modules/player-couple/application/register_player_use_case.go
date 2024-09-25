package application

import (
	"github.com/paguerre3/goddd/modules/player-couple/domain"
)

type RegisterPlayerUseCase interface {
	RegisterPlayerUseCase(inputPlayer domain.Player) (status RegisterPlayerStatus, err error)
}

type RegisterPlayerStatus uint8

const (
	RegisterPlayerPending RegisterPlayerStatus = iota
	RegisterPlayerUpdated
	RegisterPlayerCreated
)

func NewRegisterPlayerUseCase(playerRepository domain.PlayerRepository,
	idGenerator domain.IDGenerator) RegisterPlayerUseCase {
	return &playerService{
		playerRepo: playerRepository,
		idGen:      idGenerator,
	}
}

// RegisterPlayerUseCase registers a player or updates it if it already exists.
func (s *playerService) RegisterPlayerUseCase(inputPlayer domain.Player) (status RegisterPlayerStatus, err error) {
	// Validate new player entries.
	newPlayer, err := domain.NewPlayer(s.idGen,
		inputPlayer.Email,
		inputPlayer.SocialSecurityNumber,
		inputPlayer.FirstName,
		inputPlayer.LastName,
		inputPlayer.Age)
	if err != nil {
		return status, err
	}

	// Check if the player already exists.
	foundPlayer, err := s.findByIDOrEmail(inputPlayer.ID, inputPlayer.Email)
	if err != nil {
		return status, err
	}

	// Ensure existing player isn't an empty struct:
	if len(foundPlayer.ID) > 0 {
		err = s.playerRepo.Update(*newPlayer)
		if err == nil {
			status = RegisterPlayerUpdated
		}
	} else {
		err = s.playerRepo.Save(*newPlayer)
		if err == nil {
			status = RegisterPlayerCreated
		}
	}

	return status, err
}

// FindByIDOrEmail returns a player found by ID or email.
func (s *playerService) findByIDOrEmail(id, email string) (player domain.Player, err error) {
	if len(id) > 0 {
		if err = domain.ValidateID(id); err != nil {
			return player, err
		}
		player, err = s.playerRepo.FindByID(id)
	} else {
		// Email validation is already done at the beginning of RegisterPlayerUseCase function.
		player, err = s.playerRepo.FindByEmail(email)
	}
	return player, err
}
