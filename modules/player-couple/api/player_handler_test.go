package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paguerre3/goddd/modules/player-couple/application"
	"github.com/paguerre3/goddd/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterPlayer(t *testing.T) {
	h := &PlayerHandler{
		registerPlayerUseCase: &mockRegisterPlayerUseCase{},
	}

	tests := []struct {
		name       string
		request    string
		statusCode int
	}{
		{
			name:       "Invalid JSON binding",
			request:    `invalid json`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Invalid player data",
			request:    `{"email": "invalid"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Existing player (update)",
			request:    `{"email": "existing@example.com"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "New player (create)",
			request:    `{"email": "new@example.com"}`,
			statusCode: http.StatusCreated,
		},
		{
			name:       "Internal server error",
			request:    `{"email": "error@example.com"}`,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Invalid status",
			request:    `{"email": "invalid-status@example.com"}`,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/player", bytes.NewBuffer([]byte(test.request)))
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h.RegisterPlayer(c)

			assert.Equal(t, test.statusCode, w.Code)

			// Get the result
			result := w.Result()
			defer result.Body.Close()

			// Read the response body
			body, _ := io.ReadAll(result.Body)

			// Print the response body
			fmt.Printf("Register player handler Response: %s", body)
		})
	}
}

type mockRegisterPlayerUseCase struct{}

func (m *mockRegisterPlayerUseCase) RegisterPlayerUseCase(player domain.Player) (domain.Player, application.RegisterPlayerStatus, error) {
	switch player.Email {
	case "invalid":
		return domain.Player{}, application.RegisterPlayerInvalid, fmt.Errorf("invalid email")
	case "existing@example.com":
		return domain.Player{Email: player.Email}, application.RegisterPlayerUpdated, nil
	case "new@example.com":
		return domain.Player{Email: player.Email}, application.RegisterPlayerCreated, nil
	case "error@example.com":
		return domain.Player{}, 0, fmt.Errorf("internal server error")
	case "invalid-status@example.com":
		return domain.Player{Email: player.Email}, application.RegisterPlayerPending, nil
	default:
		return domain.Player{}, 0, nil
	}
}

func TestUnregisterPlayer(t *testing.T) {
	h := &PlayerHandler{
		unregisterPlayerUseCase: &mockUnregisterPlayerUseCase{},
	}

	t.Run("invalid player ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "playerId", Value: "invalid-id"}}
		h.UnregisterPlayer(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("non-existent player ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "playerId", Value: "non-existent-id"}}
		h.UnregisterPlayer(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("pending status", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "playerId", Value: "pending-id"}}
		h.UnregisterPlayer(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("successful unregister", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "playerId", Value: "valid-id"}}
		h.UnregisterPlayer(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "playerId", Value: "error-id"}}
		h.UnregisterPlayer(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

type mockUnregisterPlayerUseCase struct{}

func (m *mockUnregisterPlayerUseCase) UnregisterPlayerUseCase(playerId string) (application.UnregisterPlayerStatus, error) {
	switch playerId {
	case "invalid-id":
		return application.UnregisterPlayerInvalid, errors.New("invalid player ID")
	case "non-existent-id":
		return application.UnregisterPlayerNotFound, nil
	case "pending-id":
		return application.UnregisterPlayerPending, nil
	case "valid-id":
		return application.UnregisterPlayerDeleted, nil
	case "error-id":
		return application.UnregisterPlayerDeleted, errors.New("internal server error")
	default:
		return application.UnregisterPlayerDeleted, nil
	}
}

type mockFindPlayerUseCase struct {
	mock.Mock
}

func (m *mockFindPlayerUseCase) FindPlayerByIDUseCase(playerId string) (domain.Player, application.FindPlayerStatus, error) {
	args := m.Called(playerId)
	return args.Get(0).(domain.Player), args.Get(1).(application.FindPlayerStatus), args.Error(2)
}

func (m *mockFindPlayerUseCase) FindPlayerByEmailUseCase(email string) (domain.Player, application.FindPlayerStatus, error) {
	args := m.Called(email)
	return args.Get(0).(domain.Player), args.Get(1).(application.FindPlayerStatus), args.Error(2)
}

func (m *mockFindPlayerUseCase) FindPlayersByLastNameUseCase(lastName string) ([]domain.Player, application.FindPlayerStatus, error) {
	args := m.Called(lastName)
	return args.Get(0).([]domain.Player), args.Get(1).(application.FindPlayerStatus), args.Error(2)
}

func TestFindPlayerByID(t *testing.T) {
	// Arrange
	findPlayerUseCaseMock := &mockFindPlayerUseCase{}
	playerHandler := &PlayerHandler{
		findPlayerUseCase: findPlayerUseCaseMock,
	}

	// Valid player ID
	t.Run("Valid player ID", func(t *testing.T) {
		// Arrange
		playerId := "valid-id"
		foundPlayer := domain.Player{ID: playerId}
		findPlayerUseCaseMock.On("FindPlayerByIDUseCase", playerId).Return(foundPlayer, application.FindPlayerFound, nil)

		// Act
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{
			{Key: "playerId", Value: playerId},
		}
		playerHandler.FindPlayerByID(c)

		// Assert
		assert.Equal(t, http.StatusOK, c.Writer.Status())
	})

	// Invalid player ID
	t.Run("Invalid player ID", func(t *testing.T) {
		// Arrange
		playerId := "invalid-id"
		findPlayerUseCaseMock.On("FindPlayerByIDUseCase", playerId).Return(domain.Player{}, application.FindPlayerInvalid, errors.New("invalid ID"))

		// Act
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{
			{Key: "playerId", Value: playerId},
		}
		playerHandler.FindPlayerByID(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})

	// Player not found
	t.Run("Player not found", func(t *testing.T) {
		// Arrange
		playerId := "not-found-id"
		findPlayerUseCaseMock.On("FindPlayerByIDUseCase", playerId).Return(domain.Player{}, application.FindPlayerNotFound, nil)

		// Act
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{
			{Key: "playerId", Value: playerId},
		}
		playerHandler.FindPlayerByID(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, c.Writer.Status())
	})

	// Error in finding player
	t.Run("Error in finding player", func(t *testing.T) {
		// Arrange
		playerId := "error-id"
		findPlayerUseCaseMock.On("FindPlayerByIDUseCase", playerId).Return(domain.Player{}, application.FindPlayerPending, errors.New("error in finding player"))

		// Act
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Params = gin.Params{
			{Key: "playerId", Value: playerId},
		}
		playerHandler.FindPlayerByID(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	})
}
