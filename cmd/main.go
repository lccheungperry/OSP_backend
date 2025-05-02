package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lccheungperry/OSP_backend/internal/infrastructure/mongodb/connection"
	"github.com/lccheungperry/OSP_backend/internal/interfaces/http/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	mongoClient, err := connection.NewMongoClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Get database name from environment or use default
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "survey_platform"
	}

	// Get database
	db := mongoClient.Database(dbName)

	router := gin.Default()

	handlers.SetupRoutes(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
