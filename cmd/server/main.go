package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohammadbafkar/web-analyzer/internal/handlers"
)

func main() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.GET("/healthz", handlers.HandleHealth)
	r.GET("/readyz", handlers.HandleReady)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
