package application

import (
	"github.com/paguerre3/goddd/internal/modules/player-couple/domain"
)

type RegisterPlayerUseCase interface {
	RegisterPlayerUseCase(inputPlayer domain.Player) (newPlayer domain.Player, status RegisterPlayerStatus, err error)
}

type RegisterPlayerStatus uint8

const (
	RegisterPlayerPending RegisterPlayerStatus = iota
	RegisterPlayerInvalid
	RegisterPlayerUpdated
	RegisterPlayerCreated
)

func NewRegisterPlayerUseCase(idGenerator domain.IDGenerator,
	playerRepository domain.PlayerRepository) RegisterPlayerUseCase {
	return &playerService{
		idGen:      idGenerator,
		playerRepo: playerRepository,
	}
}

// RegisterPlayerUseCase registers a player or updates it if it already exists.
func (s *playerService) RegisterPlayerUseCase(inputPlayer domain.Player) (newPlayer domain.Player,
	status RegisterPlayerStatus, err error) {
	// Validate new player entries.
	newPlayerRef, err := domain.NewPlayer(s.idGen,
		inputPlayer.Email,
		inputPlayer.SocialSecurityNumber,
		inputPlayer.FirstName,
		inputPlayer.LastName,
		inputPlayer.Age)
	if err != nil {
		status = RegisterPlayerInvalid
		return newPlayer, status, err
	}

	// Check if the player already exists.
	foundPlayer, status, err := s.findByIDOrEmail(inputPlayer.ID, inputPlayer.Email)
	if err != nil {
		return newPlayer, status, err
	}

	// Ensure existing player isn't an empty struct:
	if len(foundPlayer.ID) > 0 {
		// Ensure to overwrite auto generated ID of new player.
		newPlayerRef.ID = foundPlayer.ID
		err = s.playerRepo.Update(*newPlayerRef)
		if err == nil {
			status = RegisterPlayerUpdated
		}
	} else {
		// A valid ID never overwrites the auto generated one during creation.
		err = s.playerRepo.Save(*newPlayerRef)
		if err == nil {
			status = RegisterPlayerCreated
		}
	}
	if err == nil {
		newPlayer = *newPlayerRef
	}
	return newPlayer, status, err
}

// FindByIDOrEmail returns a player found by ID or email.
func (s *playerService) findByIDOrEmail(id, email string) (player domain.Player,
	status RegisterPlayerStatus, err error) {
	if len(id) > 0 {
		if err = domain.ValidateID(id); err != nil {
			status = RegisterPlayerInvalid
			return player, status, err
		}
		player, err = s.playerRepo.FindByID(id)
	} else {
		// Email validation is already done at the beginning of RegisterPlayerUseCase function.
		player, err = s.playerRepo.FindByEmail(email)
	}
	return player, status, err
}
