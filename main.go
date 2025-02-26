package main

import (
	"super-fantasy-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// API Versioning
	v1 := router.Group("/api/v1")
	{
		// Basketball routes
		basketball := v1.Group("/basketball")
		basketball.POST("/projections", handlers.CalculateBasketballProjections)

		// Hockey routes
		hockey := v1.Group("/hockey")
		hockey.POST("/projections", handlers.CalculateHockeyProjections)

		// Football routes
		football := v1.Group("/football")
		football.POST("/projections", handlers.CalculateFootballProjections)

		// Baseball routes
		baseball := v1.Group("/baseball")
		baseball.POST("/projections", handlers.CalculateBaseballProjections)
	}

	router.Run(":8080")
}