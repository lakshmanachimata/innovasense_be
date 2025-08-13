package controllers

import (
	"innovasense_be/models"
	"innovasense_be/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HistoricalDataController struct {
	historicalDataService *services.HistoricalDataService
	orgService            *services.OrganizationService
}

func NewHistoricalDataController(
	historicalDataService *services.HistoricalDataService,
	orgService *services.OrganizationService,
) *HistoricalDataController {
	return &HistoricalDataController{
		historicalDataService: historicalDataService,
		orgService:            orgService,
	}
}

// GetHistoricalData handles the historical data request
// @Summary Get historical data
// @Description Get historical data for a user based on organization credentials
// @Tags Historical Data
// @Accept json
// @Produce json
// @Param apikey header string true "API Key"
// @Param secretkey header string true "Secret Key"
// @Param request body models.HistoricalDataRequest true "Historical data request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /Services/getHistoricalData [post]
func (c *HistoricalDataController) GetHistoricalData(ctx *gin.Context) {
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
	var req models.HistoricalDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Step 1: Validate organization credentials
	org, err := c.orgService.ValidateOrgCredentials(apiKey, secretKey)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "Organization validation failed: " + err.Error(),
		})
		return
	}

	// Step 2: Check if user exists for this organization
	orgUser, err := c.orgService.CheckUserExistsByContact(req.Contact, org.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    1,
			Message: "Failed to check user existence: " + err.Error(),
		})
		return
	}

	if orgUser == nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Code:    1,
			Message: "User not found",
		})
		return
	}

	// Step 3: Get user ID from users_master table
	userID, err := c.orgService.GetUserIDByContact(req.Contact)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    1,
			Message: "Failed to get user ID: " + err.Error(),
		})
		return
	}

	// Step 4: Get historical data
	historicalData, err := c.historicalDataService.GetHistoricalData(userID, req.FromDate, req.ToDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    1,
			Message: "Failed to retrieve historical data: " + err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "Success",
		Response: historicalData,
	})
}
