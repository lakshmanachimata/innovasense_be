package controllers

import (
	"innovasense_be/models"
	"innovasense_be/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HydrationRecommendationController struct {
	hydrationRecommendationService *services.HydrationRecommendationService
}

func NewHydrationRecommendationController(
	hydrationRecommendationService *services.HydrationRecommendationService,
) *HydrationRecommendationController {
	return &HydrationRecommendationController{
		hydrationRecommendationService: hydrationRecommendationService,
	}
}

// GetHydrationRecommendation handles the hydration recommendation request
// @Summary Get hydration recommendation
// @Description Get hydration recommendation based on user data and organization credentials
// @Tags Hydration
// @Accept json
// @Produce json
// @Param apikey header string true "API Key"
// @Param secretkey header string true "Secret Key"
// @Param request body models.HydrationRecommendationRequest true "Hydration recommendation request"
// @Success 200 {object} models.EnhancedHydrationResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /Services/getHydrationRecommendation [post]
func (c *HydrationRecommendationController) GetHydrationRecommendation(ctx *gin.Context) {
	// Extract API key and secret key from headers
	apiKey := ctx.GetHeader("apikey")
	secretKey := ctx.GetHeader("secretkey")

	if apiKey == "" || secretKey == "" {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "API key and secret key are required in headers",
		})
		return
	}

	// Parse request body
	var req models.HydrationRecommendationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Process the request
	response, err := c.hydrationRecommendationService.GetHydrationRecommendation(&req, apiKey, secretKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    1,
			Message: "Failed to process request: " + err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "Success",
		Response: response,
	})
}
