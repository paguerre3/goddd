package application

import (
	"errors"
	"testing"

	"github.com/paguerre3/goddd/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
)

func TestFindPlayerByIDUseCase(t *testing.T) {
	t.Run("Valid player ID found", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewFindPlayerUseCase(repo, idGen)
		playerId := "valid-id"
		foundPlayer := domain.Player{ID: playerId}
		repo.On("FindByID", playerId).Return(foundPlayer, nil)

		// Act
		player, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, foundPlayer, player)
	})

	t.Run("Invalid player ID", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewFindPlayerUseCase(repo, idGen)
		playerId := "i" // invalid id
		expectedErr := domain.ValidateID(playerId)

		// Act
		player, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.Player{}, player)
	})

	t.Run("Player not found", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewFindPlayerUseCase(repo, idGen)
		playerId := "not-found-id"
		repo.On("FindByID", playerId).Return(domain.Player{}, nil)

		// Act
		player, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, domain.Player{}, player)
	})

	t.Run("Error in repository", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewFindPlayerUseCase(repo, idGen)
		playerId := "error-id"
		expectedErr := errors.New("repo error")
		repo.On("FindByID", playerId).Return(domain.Player{}, expectedErr)

		// Act
		player, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, domain.Player{}, player)
	})
}
