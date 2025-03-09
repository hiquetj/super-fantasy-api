package main

import (
	"log"
	"os"
	"super-fantasy-api/db"
	"super-fantasy-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults or environment variables")
	}

	// Fetch MongoDB configuration from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	// Initialize MongoDB with environment variables
	if err := db.InitMongoDB(mongoURI, dbName, collectionName); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}
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
		uploadBaseball := baseball.Group("/upload")
		uploadBaseball.POST("", handlers.UploadCSV)
	}

	router.Run(":8080")
}
