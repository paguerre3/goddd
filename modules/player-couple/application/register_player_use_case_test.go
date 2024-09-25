package application

import (
	"fmt"
	"testing"

	"github.com/paguerre3/goddd/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mockId = "mock-id"
)

type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Save(player domain.Player) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *MockPlayerRepository) FindByID(id string) (domain.Player, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Player), args.Error(1)
}

func (m *MockPlayerRepository) FindByEmail(email string) (domain.Player, error) {
	args := m.Called(email)
	return args.Get(0).(domain.Player), args.Error(1)
}

func (m *MockPlayerRepository) FindByLastName(lastName string) ([]domain.Player, error) {
	args := m.Called(lastName)
	return args.Get(0).([]domain.Player), args.Error(1)
}

func (m *MockPlayerRepository) Update(player domain.Player) error {
	args := m.Called(player)
	return args.Error(0)
}

func (m *MockPlayerRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock IDGenerator for testing without breaking modularity principle even in testing domain module package.
type MockIDGenerator struct {
	// molck aggregate used as a helper mock to avoid repetitions:
	aggregate *string
}

func (m *MockIDGenerator) GenerateID() string {
	if m.aggregate != nil {
		return fmt.Sprintf("%s-%s", mockId, *m.aggregate)
	}
	return mockId
}

func (m *MockIDGenerator) GenerateIDWithPrefixes(p1, p2 string) string {
	return fmt.Sprintf("%s-%s-%s", p1, p2, mockId)
}

func TestUpsertPlayerUseCase_Success(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{}, nil)
	repo.On("Save", mock.Anything).Return(nil)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.NoError(t, err)
}

func TestUpsertPlayerUseCase_UpdateExistingPlayerByID(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{
		ID:        "existing-id",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByID", inputPlayer.ID).Return(domain.Player{ID: inputPlayer.ID}, nil)
	repo.On("Update", mock.Anything).Return(nil)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.NoError(t, err)
}

func TestUpsertPlayerUseCase_UpdateExistingPlayerByEmail(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{ID: "existing-id"}, nil)
	repo.On("Update", mock.Anything).Return(nil)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.NoError(t, err)
}

func TestUpsertPlayerUseCase_ValidationError(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{Email: ""}

	// Expect
	_, expectedErr := domain.NewPlayer(idGen, inputPlayer.Email, nil, "", "", nil)
	assert.Error(t, expectedErr)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUpsertPlayerUseCase_FindByIDError(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{
		ID:        "existing-id",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	expectedErr := assert.AnError
	repo.On("FindByID", inputPlayer.ID).Return(domain.Player{}, expectedErr)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUpsertPlayerUseCase_FindByEmailError(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	expectedErr := assert.AnError
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{}, expectedErr)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUpsertPlayerUseCase_SaveError(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByEmail", inputPlayer.Email).Return(domain.Player{}, nil)
	expectedErr := assert.AnError
	repo.On("Save", mock.Anything).Return(expectedErr)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestUpsertPlayerUseCase_UpdateError(t *testing.T) {
	// Arrange
	idGen := &MockIDGenerator{}
	repo := &MockPlayerRepository{}
	service := &PlayerCoupleService{idGen: idGen, playerRepo: repo}
	inputPlayer := domain.Player{
		ID:        "existing-id",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
	}

	// Expect
	repo.On("FindByID", inputPlayer.ID).Return(domain.Player{ID: inputPlayer.ID}, nil)
	expectedErr := assert.AnError
	repo.On("Update", mock.Anything).Return(expectedErr)

	// Act
	err := service.UpsertPlayerUseCase(inputPlayer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
