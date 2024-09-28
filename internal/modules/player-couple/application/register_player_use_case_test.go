package application

import (
	"fmt"
	"testing"

	"github.com/paguerre3/goddd/internal/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mockId = "mock-id"
)

type mockPlayerRepository struct {
	mock.Mock
}

func (m *mockPlayerRepository) Upsert(player *domain.Player) error {
	args := m.Called(player)
	if player.ID == "" {
		idGen := mockIDGenerator{}
		player.ID = idGen.GenerateID()
	}
	return args.Error(0)
}

func (m *mockPlayerRepository) FindByID(id string) (domain.Player, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Player), args.Error(1)
}

func (m *mockPlayerRepository) FindByEmail(email string) (domain.Player, error) {
	args := m.Called(email)
	return args.Get(0).(domain.Player), args.Error(1)
}

func (m *mockPlayerRepository) FindByLastName(lastName string) ([]domain.Player, error) {
	args := m.Called(lastName)
	return args.Get(0).([]domain.Player), args.Error(1)
}

func (m *mockPlayerRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock IDGenerator for testing without breaking modularity principle even in testing domain module package.
type mockIDGenerator struct {
	// molck aggregate used as a helper mock to avoid repetitions:
	aggregate *string
}

func (m *mockIDGenerator) GenerateID() string {
	if m.aggregate != nil {
		return fmt.Sprintf("%s-%s", mockId, *m.aggregate)
	}
	return mockId
}

func (m *mockIDGenerator) GenerateIDWithPrefixes(p1, p2 string) string {
	return fmt.Sprintf("%s-%s-%s", p1, p2, mockId)
}

func TestRegisterPlayerUseCase_Success(t *testing.T) {
	// Arrange
	idGen := &mockIDGenerator{}
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{}, nil)
	repo.On("Upsert", mock.Anything).Return(nil)
	expectedNewPlayer := domain.Player{
		ID:        idGen.GenerateID(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, RegisterPlayerCreated, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestRegisterPlayerUseCase_UpdateExistingPlayerByID(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		ID:        "existing-id",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByID", inputPlayer.ID).Return(domain.Player{ID: inputPlayer.ID}, nil)
	repo.On("Upsert", mock.Anything).Return(nil)
	expectedNewPlayer := inputPlayer

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, RegisterPlayerUpdated, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestRegisterPlayerUseCase_UpdateExistingPlayerByIDValidationError(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		// Invalid ID
		ID:        "i",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	expectedErr := domain.ValidateID(inputPlayer.ID)
	var expectedNewPlayer domain.Player

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, RegisterPlayerInvalid, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestRegisterPlayerUseCase_UpdateExistingPlayerByEmail(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{ID: "existing-id"}, nil)
	repo.On("Upsert", mock.Anything).Return(nil)
	expectedNewPlayer := domain.Player{
		ID:        "existing-id",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, RegisterPlayerUpdated, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestRegisterPlayerUseCase_ValidationError(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{Email: ""}

	// Expect
	_, expectedErr := domain.NewPlayer(inputPlayer.Email, nil, "", "", nil)
	assert.Error(t, expectedErr)
	var expectedNewPlayer domain.Player

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, RegisterPlayerInvalid, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestRegisterPlayerUseCase_FindByIDError(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		ID:        "existing-id",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	expectedErr := assert.AnError
	repo.On("FindByID", inputPlayer.ID).Return(domain.Player{}, expectedErr)
	var expectedNewPlayer domain.Player

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, RegisterPlayerPending, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestRegisterPlayerUseCase_FindByEmailError(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	expectedErr := assert.AnError
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{}, expectedErr)
	var expectedNewPlayer domain.Player

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, RegisterPlayerPending, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestRegisterPlayerUseCase_SaveError(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{}, nil)
	expectedErr := assert.AnError
	repo.On("Upsert", mock.Anything).Return(expectedErr)
	var expectedNewPlayer domain.Player

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, RegisterPlayerPending, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}

func TestPlayerUseCase_UpdateError(t *testing.T) {
	// Arrange
	repo := &mockPlayerRepository{}
	service := NewRegisterPlayerUseCase(repo)
	inputPlayer := domain.Player{
		ID:        "existing-id",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByID", inputPlayer.ID).Return(domain.Player{ID: inputPlayer.ID}, nil)
	expectedErr := assert.AnError
	repo.On("Upsert", mock.Anything).Return(expectedErr)
	var expectedNewPlayer domain.Player

	// Act
	newPlayer, status, err := service.RegisterPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, RegisterPlayerPending, status)
	assert.Equal(t, expectedNewPlayer, newPlayer)
}
