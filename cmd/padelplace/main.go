package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paguerre3/goddd/internal/modules/common/mongo"
	"github.com/paguerre3/goddd/internal/modules/common/utils"
	"github.com/paguerre3/goddd/internal/modules/player-couple/api"
	"github.com/paguerre3/goddd/internal/modules/player-couple/application"
	player_couple_infrastructure "github.com/paguerre3/goddd/internal/modules/player-couple/infrastructure/mongo"
)

func main() {
	mongoClient := mongo.NewMongoClient()
	defer mongoClient.Close()

	idGen := utils.NewUUIDGenerator()

	playerRepo := player_couple_infrastructure.NewMongoPlayerRepository(mongoClient)

	registerPlayerUseCase := application.NewRegisterPlayerUseCase(playerRepo, idGen)
	unregisterPlayerUseCase := application.NewUnregisterPlayerUseCase(playerRepo, idGen)
	findPlayerUseCase := application.NewFindPlayerUseCase(playerRepo, idGen)

	playerHandler := api.NewPlayerHandler(registerPlayerUseCase, unregisterPlayerUseCase, findPlayerUseCase)

	// Initialize router
	router := gin.Default()

	// Routes
	router.POST("/players", playerHandler.RegisterPlayer)
	router.DELETE("/players/:playerId", playerHandler.UnregisterPlayer)
	router.GET("/players/:playerId", playerHandler.FindPlayerByID)
	router.GET("/players/email/:email", playerHandler.FindPlayerByEmail)
	router.GET("/players/last-name/:lastName", playerHandler.FindPlayersByLastName)

	// Start your HTTP server and handle routes
	router.Run(":8080")
}
