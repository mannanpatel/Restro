package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mannanpatel/Restro/common/database"
	"github.com/mannanpatel/Restro/config"
	"github.com/mannanpatel/Restro/internal/routes"
)

func main() {
	// env load
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// config load
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// db connect
	database.DBConnection(cfg)

	// gin mode
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.SetupRoutes(r)

	// Step 6 — Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Restro API is healthy!"})
	})

	r.Run(":" + cfg.Server.Port) // listen and serve on 0.0.0.0:8085 (for windows "localhost:8085")
}
