package application

import (
	"errors"
	"testing"

	"github.com/paguerre3/goddd/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
)

func TestUnregisterPlayerUseCase(t *testing.T) {
	t.Run("Valid player ID found unregistered successfully", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "valid-id"
		foundPlayer := domain.Player{ID: playerId}
		repo.On("FindByID", playerId).Return(foundPlayer, nil)
		repo.On("Delete", playerId).Return(nil)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, UnregisterPlayerDeleted, status)
	})

	t.Run("Invalid player ID", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "i"
		expectedErr := domain.ValidateID(playerId)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, UnregisterPlayerInvalid, status)
	})

	t.Run("Player not found", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "not-found-id"
		repo.On("FindByID", playerId).Return(domain.Player{}, nil)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, UnregisterPlayerNotFound, status)
	})

	t.Run("Error finding player by ID", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "error-id"
		expectedErr := errors.New("error finding player")
		repo.On("FindByID", playerId).Return(domain.Player{}, expectedErr)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, UnregisterPlayerPending, status)
	})

	t.Run("Error deleting player", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "delete-error-id"
		foundPlayer := domain.Player{ID: playerId}
		repo.On("FindByID", playerId).Return(foundPlayer, nil)
		expectedErr := errors.New("error deleting player")
		repo.On("Delete", playerId).Return(expectedErr)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, UnregisterPlayerPending, status)
	})
}
