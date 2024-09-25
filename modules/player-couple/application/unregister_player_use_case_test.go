package application

import (
	"testing"

	"github.com/paguerre3/goddd/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
)

func TestUnregisterPlayerUseCase(t *testing.T) {
	t.Run("Existing player ID", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "existing-id"
		foundPlayer := domain.Player{ID: playerId}

		repo.On("FindByID", playerId).Return(foundPlayer, nil)
		repo.On("Delete", playerId).Return(nil)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, UnregisterPlayerDeleted, status)
	})

	t.Run("Non-existent player ID", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "non-existent-id"
		foundPlayer := domain.Player{}

		repo.On("FindByID", playerId).Return(foundPlayer, nil)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, UnregisterPlayerNotFound, status)
	})

	t.Run("Error in FindByID", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "existing-id"
		expectedErr := assert.AnError

		repo.On("FindByID", playerId).Return(domain.Player{}, expectedErr)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, UnregisterPlayerPending, status)
	})

	t.Run("Error in Delete", func(t *testing.T) {
		// Arrange
		idGen := &MockIDGenerator{}
		repo := &MockPlayerRepository{}
		service := NewUnregisterPlayerUseCase(repo, idGen)
		playerId := "existing-id"
		foundPlayer := domain.Player{ID: playerId}
		expectedErr := assert.AnError

		repo.On("FindByID", playerId).Return(foundPlayer, nil)
		repo.On("Delete", playerId).Return(expectedErr)

		// Act
		status, err := service.UnregisterPlayerUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, UnregisterPlayerPending, status)
	})
}
