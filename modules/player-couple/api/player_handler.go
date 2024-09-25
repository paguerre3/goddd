package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paguerre3/goddd/modules/player-couple/application"
	"github.com/paguerre3/goddd/modules/player-couple/domain"
)

type PlayerHandler struct {
	registerPlayerUseCase   application.RegisterPlayerUseCase
	unregisterPlayerUseCase application.UnregisterPlayerUseCase
	findPlayerUseCase       application.FindPlayerUseCase
}

func NewPlayerHandler(registerPlayerUseCase application.RegisterPlayerUseCase,
	unregisterPlayerUseCase application.UnregisterPlayerUseCase,
	findPlayerUseCase application.FindPlayerUseCase) *PlayerHandler {
	return &PlayerHandler{
		registerPlayerUseCase:   registerPlayerUseCase,
		unregisterPlayerUseCase: unregisterPlayerUseCase,
		findPlayerUseCase:       findPlayerUseCase,
	}
}

func (h *PlayerHandler) RegisterPlayer(c *gin.Context) {
	var player domain.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newPlayer, status, err := h.registerPlayerUseCase.RegisterPlayerUseCase(player)
	if err != nil {
		errCode := http.StatusInternalServerError
		if status == application.RegisterPlayerInvalid {
			errCode = http.StatusBadRequest
		}
		c.JSON(errCode, gin.H{"error": err.Error()})
		return
	}
	switch status {
	case application.RegisterPlayerUpdated:
		c.JSON(http.StatusOK, newPlayer)
	case application.RegisterPlayerCreated:
		c.JSON(http.StatusCreated, newPlayer)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("invalid status %d", status)})
	}
}

func (h *PlayerHandler) UnregisterPlayer(c *gin.Context) {
	playerId := c.Param("playerId")
	status, err := h.unregisterPlayerUseCase.UnregisterPlayerUseCase(playerId)
	if err != nil {
		errCode := http.StatusInternalServerError
		if status == application.UnregisterPlayerInvalid {
			errCode = http.StatusBadRequest
		}
		c.JSON(errCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}
