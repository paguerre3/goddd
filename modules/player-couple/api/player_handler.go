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
		if status == application.RegisterPlayerInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		if status == application.UnregisterPlayerInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if status == application.UnregisterPlayerNotFound {
		c.JSON(http.StatusNotFound, gin.H{"status": status.String()})
		return
	}
	if status == application.UnregisterPlayerPending {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("invalid status %d", status)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status.String()})
}

func (h *PlayerHandler) FindPlayerByID(c *gin.Context) {
	playerId := c.Param("playerId")
	player, status, err := h.findPlayerUseCase.FindPlayerByIDUseCase(playerId)
	handleFindResponse(c, player, status, err)
}

func (h *PlayerHandler) FindPlayerByEmail(c *gin.Context) {
	email := c.Param("email")
	player, status, err := h.findPlayerUseCase.FindPlayerByEmailUseCase(email)
	handleFindResponse(c, player, status, err)
}

func (h *PlayerHandler) FindPlayersByLastName(c *gin.Context) {
	lastName := c.Param("lastName")
	players, status, err := h.findPlayerUseCase.FindPlayersByLastNameUseCase(lastName)
	handleFindResponse(c, players, status, err)
}

func handleFindResponse[T domain.Player | []domain.Player](c *gin.Context, playerS T, status application.FindPlayerStatus, err error) {
	if err != nil {
		if status == application.FindPlayerInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch status {
	case application.FindPlayerNotFound:
		c.JSON(http.StatusNotFound, playerS)
	case application.FindPlayerFound:
		c.JSON(http.StatusOK, playerS)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("invalid status %d", status)})
	}
}
