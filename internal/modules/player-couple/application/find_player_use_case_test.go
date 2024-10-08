package application

import (
	"errors"
	"testing"

	"github.com/paguerre3/goddd/internal/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
)

func TestFindPlayerByIDUseCase(t *testing.T) {
	t.Run("Valid ID player found", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		playerId := "valid-id"
		foundPlayer := domain.Player{ID: playerId}
		repo.On("FindByID", playerId).Return(foundPlayer, nil)

		// Act
		player, status, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, FindPlayerFound, status)
		assert.Equal(t, foundPlayer, player)
	})

	t.Run("Invalid player ID", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		playerId := "i" // invalid id
		expectedErr := domain.ValidateID(playerId)

		// Act
		player, status, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, FindPlayerInvalid, status)
		assert.Equal(t, domain.Player{}, player)
	})

	t.Run("Player not found", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		playerId := "not-found-id"
		repo.On("FindByID", playerId).Return(domain.Player{}, nil)

		// Act
		player, status, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, FindPlayerNotFound, status)
		assert.Equal(t, domain.Player{}, player)
	})

	t.Run("Error in repository finding by ID", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		playerId := "error-id"
		expectedErr := errors.New("repo error")
		repo.On("FindByID", playerId).Return(domain.Player{}, expectedErr)

		// Act
		player, status, err := service.FindPlayerByIDUseCase(playerId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, FindPlayerPending, status)
		assert.Equal(t, domain.Player{}, player)
	})
}

func TestFindPlayerByEmailUseCase(t *testing.T) {
	t.Run("Valid email player found", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		email := "test@example.com"
		foundPlayer := domain.Player{ID: "1234567", Email: email}
		repo.On("FindByEmail", email).Return(foundPlayer, nil)

		// Act
		player, status, err := service.FindPlayerByEmailUseCase(email)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, FindPlayerFound, status)
		assert.Equal(t, foundPlayer, player)
	})

	t.Run("Valid email not player found", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		email := "test@example.com"
		repo.On("FindByEmail", email).Return(domain.Player{}, nil)

		// Act
		player, status, err := service.FindPlayerByEmailUseCase(email)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, FindPlayerNotFound, status)
		assert.Equal(t, domain.Player{}, player)
	})

	t.Run("Invalid email", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		email := "invalid-email"
		expectedErr := domain.ValidateEmail(email)

		// Act
		player, status, err := service.FindPlayerByEmailUseCase(email)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, FindPlayerInvalid, status)
		assert.Equal(t, domain.Player{}, player)
	})

	t.Run("Error in repository finding by email", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		email := "test@example.com"
		repoErr := errors.New("repository error")
		repo.On("FindByEmail", email).Return(domain.Player{}, repoErr)

		// Act
		player, status, err := service.FindPlayerByEmailUseCase(email)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Equal(t, FindPlayerPending, status)
		assert.Equal(t, domain.Player{}, player)
	})
}

func TestFindPlayersByLastNameUseCase(t *testing.T) {
	t.Run("Valid last name players found", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		lastName := "Doe"
		expectedPlayers := []domain.Player{{ID: "1234567", LastName: "Doe"}}

		repo.On("FindByLastName", lastName).Return(expectedPlayers, nil)

		// Act
		players, status, err := service.FindPlayersByLastNameUseCase(lastName)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, FindPlayerFound, status)
		assert.Equal(t, expectedPlayers, players)
	})

	t.Run("Valid last name players not found", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		lastName := "Doe"
		var expectedPlayers []domain.Player = nil

		repo.On("FindByLastName", lastName).Return(expectedPlayers, nil)

		// Act
		players, status, err := service.FindPlayersByLastNameUseCase(lastName)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, FindPlayerNotFound, status)
		assert.Equal(t, expectedPlayers, players)
	})

	t.Run("Invalid last name players not browsed", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		lastName := ""
		expectedErr := domain.ValidateLastName(lastName)

		// Act
		players, status, err := service.FindPlayersByLastNameUseCase(lastName)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, FindPlayerInvalid, status)
		assert.Nil(t, players)
	})

	t.Run("Error in repository finding by last name", func(t *testing.T) {
		// Arrange
		repo := &mockPlayerRepository{}
		service := NewFindPlayerUseCase(repo)
		lastName := "Doe"
		expectedErr := errors.New("repo error")

		// TODO: fix this
		//repo.On("FindByLastName", lastName).Return(nil, expectedErr)
		repo.On("FindByLastName", lastName).Return([]domain.Player{}, expectedErr)

		// Act
		players, status, err := service.FindPlayersByLastNameUseCase(lastName)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		// TODO: fix this
		//assert.Nil(t, players)
		assert.Equal(t, FindPlayerPending, status)
		assert.Equal(t, []domain.Player{}, players)
	})
}
