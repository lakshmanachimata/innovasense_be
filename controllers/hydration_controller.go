package controllers

import (
	"innovasense_be/middleware"
	"innovasense_be/models"
	"innovasense_be/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HydrationController struct {
	hydrationService *services.HydrationService
	commonService    *services.CommonService
	userService      *services.UserService
}

func NewHydrationController() *HydrationController {
	return &HydrationController{
		hydrationService: services.NewHydrationService(),
		commonService:    services.NewCommonService(),
		userService:      services.NewUserService(),
	}
}

// InnovoHydration handles hydration data submission
// @Summary Record hydration data
// @Description Record basic hydration data with automatic BMI, TBSA, sweat loss, and sweat rate calculations
// @Tags Hydration
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.HydrationRequest true "Hydration data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/innovoHyderation [post]
func (c *HydrationController) InnovoHydration(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.HydrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Debug: Log the received request
	log.Printf("Received request: %+v", req)
	log.Printf("JWT Claims - Email: %s, Username: %s", claims.Email, claims.UserName)

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get user ID from Email
	userID, err := c.userService.GetUserIDByEmail(claims.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	// Override UserID with the authenticated user's ID
	req.UserID = userID

	id, err := c.hydrationService.SaveHydrationData(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to save hydration data",
			Response: 0,
		})
		return
	}

	// Get the saved data to return in response
	savedData, err := c.hydrationService.GetHydrationDataByID(id)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to retrieve saved data",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "Success",
		Response: savedData,
	})
}

// NewInnovoHydration handles enhanced hydration data submission
// @Summary Record enhanced hydration data
// @Description Record enhanced hydration data with additional calculations and summaries
// @Tags Hydration
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.HydrationRequest true "Enhanced hydration data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/newinnovoHyderation [post]
func (c *HydrationController) NewInnovoHydration(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.HydrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Debug: Log the received request
	log.Printf("Received request: %+v", req)
	log.Printf("JWT Claims - Email: %s, Username: %s", claims.Email, claims.UserName)

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get user ID from Email
	userID, err := c.userService.GetUserIDByEmail(claims.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	// Override UserID with the authenticated user's ID
	req.UserID = userID

	enhancedResponse, err := c.hydrationService.SaveEnhancedHydrationData(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to save hydration data",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:             0,
		Message:          "Success",
		Response:         enhancedResponse.Data,
		SweatSummary:     enhancedResponse.SweatSummary,
		SweatRateSummary: enhancedResponse.SweatRateSummary,
	})
}

// UpdateHydrationValue handles hydration data updates
// @Summary Update hydration data
// @Description Update existing hydration data
// @Tags Hydration
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.UpdateHydrationRequest true "Updated hydration data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/updateHyderationValue [post]
func (c *HydrationController) UpdateHydrationValue(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.UpdateHydrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get user ID from Email
	userID, err := c.userService.GetUserIDByEmail(claims.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	// Verify that the user owns this hydration record
	// You might want to add additional validation here to ensure the user can only update their own data
	// For now, we'll just log the authenticated user ID for debugging
	_ = userID // Use the variable to avoid linter error

	err = c.hydrationService.UpdateHydrationData(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to update hydration data",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "Hydration data updated successfully",
	})
}

// UpdateSweatData handles sweat data updates
// @Summary Update sweat data
// @Description Update sweat analysis data
// @Tags Hydration
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.UpdateSweatDataRequest true "Sweat data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/updateSweatData [post]
func (c *HydrationController) UpdateSweatData(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.UpdateSweatDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get user ID from Email
	userID, err := c.userService.GetUserIDByEmail(claims.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	// Override UserID with the authenticated user's ID
	req.UserID = userID

	err = c.hydrationService.UpdateSweatData(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to update sweat data",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "Sweat data updated successfully",
	})
}

// GetSummary handles summary retrieval
// @Summary Get summary
// @Description Get summary data based on sweat position
// @Tags Reports
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.SummaryRequest true "Summary request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getSummary [post]
func (c *HydrationController) GetSummary(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.SummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	summary, err := c.hydrationService.GetSummary(req.SweatPosition)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get summary",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: summary,
	})
}

// GetUserDetailedSummary handles detailed summary retrieval (matches PHP logic with multiple calls)
// @Summary Get detailed summary
// @Description Get detailed summary with multiple data sources (hydration data, summary, sweat summary, sweat rate summary)
// @Tags Reports
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.DetailedSummaryRequest true "Detailed summary request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getUserDetailedSummary [post]
func (c *HydrationController) GetUserDetailedSummary(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.DetailedSummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	detailedSummary, err := c.hydrationService.GetUserDetailedSummary(req.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get detailed summary",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:             0,
		Message:          "Success",
		Response:         detailedSummary.Data,
		Summary:          detailedSummary.Summary,
		SweatSummary:     detailedSummary.SweatSummary,
		SweatRateSummary: detailedSummary.SweatRateSummary,
	})
}

// GetHydrationSummaryScreen retrieves formatted data for the summary screen
// @Summary Get hydration summary screen data
// @Description Get formatted hydration data specifically for the summary screen display
// @Tags Hydration
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.HydrationSummaryRequest true "Summary screen request data"
// @Success 200 {object} models.HydrationSummaryResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getHydrationSummaryScreen [post]
func (c *HydrationController) GetHydrationSummaryScreen(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.HydrationSummaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get the summary screen data
	summaryData, err := c.hydrationService.GetHydrationSummaryScreen(req.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get summary screen data",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.HydrationSummaryResponse{
		Code:     0,
		Message:  "Success",
		Response: summaryData,
	})
}

// GetClientHistory handles client history retrieval
// @Summary Get client history
// @Description Get recent history for a client
// @Tags Reports
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.ClientHistoryRequest true "Client history request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getClientHistory [post]
func (c *HydrationController) GetClientHistory(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.ClientHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get user ID from Email
	userID, err := c.userService.GetUserIDByEmail(claims.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	history, err := c.hydrationService.GetClientHistory(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get client history",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: history,
	})
}

// GetHydrationHistory handles hydration history retrieval (matches PHP logic)
// @Summary Get hydration history
// @Description Get hydration history for a date range with inclusive date handling
// @Tags Reports
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.HistoryRequest true "History request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getHyderartionHistory [post]
func (c *HydrationController) GetHydrationHistory(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.HistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get user ID from Email
	userID, err := c.userService.GetUserIDByEmail(claims.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	// Override UserID with the authenticated user's ID
	req.UserID = userID

	history, err := c.hydrationService.GetHydrationHistory(req.UserID, req.FromDate, req.ToDate)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get hydration history",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: history,
	})
}

// GetElectrolyteHistory handles electrolyte history retrieval (matches PHP logic - only creation_datetime and image_id)
// @Summary Get electrolyte history
// @Description Get electrolyte history for a date range (returns only creation_datetime and image_id)
// @Tags Reports
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.HistoryRequest true "History request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getElectrolyteHistory [post]
func (c *HydrationController) GetElectrolyteHistory(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.HistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Decrypt the email from request body to compare with JWT claims
	decryptedEmail, err := c.userService.GetEncryptDecryptService().GetDecryptData(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email",
		})
		return
	}

	// Validate decrypted email against JWT claims
	if decryptedEmail != claims.Email {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "email in request body does not match authenticated user",
		})
		return
	}

	// Username is not encrypted, validate directly
	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	// Get user ID from Email
	userID, err := c.userService.GetUserIDByEmail(claims.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	// Override UserID with the authenticated user's ID
	req.UserID = userID

	history, err := c.hydrationService.GetElectrolyteHistory(req.UserID, req.FromDate, req.ToDate)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get electrolyte history",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: history,
	})
}
