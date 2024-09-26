package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paguerre3/goddd/modules/player-couple/application"
	"github.com/paguerre3/goddd/modules/player-couple/domain"
	"github.com/stretchr/testify/assert"
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
	default:
		return domain.Player{}, 0, nil
	}
}
