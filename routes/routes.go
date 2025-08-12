package routes

import (
	"innovasense_be/controllers"
	"innovasense_be/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Get current working directory and set up static file serving
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", currentDir)
	}

	assetsPath := currentDir + "/assets"
	log.Printf("Setting up static file serving for /assets -> %s", assetsPath)

	// Check if assets directory exists
	if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
		log.Printf("Assets directory does not exist: %s", assetsPath)
	} else {
		log.Printf("Assets directory exists: %s", assetsPath)
	}

	r.Static("/assets", assetsPath)

	// Initialize controllers
	authController := controllers.NewAuthController()
	hydrationController := controllers.NewHydrationController()
	commonController := controllers.NewCommonController()

	// API routes group
	api := r.Group("/Services")

	// Open endpoints (no authentication required)
	api.POST("/innovologin", authController.InnovoLogin)
	api.POST("/innovoregister", authController.InnovoRegister)
	api.POST("/getBannerImages", commonController.GetBannerImages)
	api.POST("/getHomeImages", commonController.GetHomeImages)

	// Protected endpoints (JWT authentication required)
	protectedGroup := api.Group("protected")
	protectedGroup.Use(middleware.JWTAuthMiddleware())
	{
		// Hydration endpoints
		protectedGroup.POST("/innovoHyderation", hydrationController.InnovoHydration)
		protectedGroup.POST("/newinnovoHyderation", hydrationController.NewInnovoHydration)
		protectedGroup.POST("/updateHyderationValue", hydrationController.UpdateHydrationValue)
		protectedGroup.POST("/updateSweatData", hydrationController.UpdateSweatData)
		protectedGroup.POST("/getSummary", hydrationController.GetSummary)
		protectedGroup.POST("/getUserDetailedSummary", hydrationController.GetUserDetailedSummary)
		protectedGroup.POST("/getHydrationSummaryScreen", hydrationController.GetHydrationSummaryScreen)
		protectedGroup.POST("/getClientHistory", hydrationController.GetClientHistory)
		protectedGroup.POST("/getHyderartionHistory", hydrationController.GetHydrationHistory)
		protectedGroup.POST("/getElectrolyteHistory", hydrationController.GetElectrolyteHistory)
		protectedGroup.POST("/uploadInnovoImage", commonController.UploadInnovoImage)
		protectedGroup.POST("/updateInnovoImagePath", commonController.UpdateInnovoImagePath)
		protectedGroup.POST("/getSweatImages", commonController.GetSweatImages)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "InnovoSens API is running",
		})
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "InnovoSens API",
			"version": "1.0.0",
		})
	})
}
