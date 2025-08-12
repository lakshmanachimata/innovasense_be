package main

import (
	"log"
	"os"

	"innovasense_be/config"
	_ "innovasense_be/docs" // This is required for swagger docs
	"innovasense_be/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           InnovoSens API
// @version         1.0
// @description     This is the InnovoSens REST API for hydration and sweat analysis data collection.
// @termsOfService  https://innovosens.com/terms

// @contact.name   InnovoSens Support
// @contact.url    https://innovosens.com
// @contact.email  support@innovosens.com

// @license.name  InnovoSens License
// @license.url   https://innovosens.com/license

// @host      localhost:8500
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// Use enhanced CORS middleware from config package
	r.Use(config.CORSMiddleware())

	// Setup routes
	routes.SetupRoutes(r)

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8500" // Changed back to 8500 as requested
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	log.Printf("API base URL: http://localhost:%s/Services", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
