package controllers

import (
	"encoding/base64"
	"fmt"
	"innovasense_be/models"
	"innovasense_be/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		userService: services.NewUserService(),
	}
}

// EncryptedLoginRequest represents the encrypted login request
type EncryptedLoginRequest struct {
	Email   string `json:"email"`
	Userpin string `json:"userpin"`
}

// Helper function to check if a string is valid base64
func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil && len(s) > 0 && !strings.Contains(s, " ")
}

// Helper function to safely decrypt or return plain text
func safeDecrypt(encryptService *services.EncryptDecryptService, data string) (string, error) {
	// If it's not base64, return as plain text
	if !isBase64(data) {
		return data, nil
	}

	// Try to decrypt, if it fails, return as plain text
	decrypted, err := encryptService.GetDecryptData(data)
	if err != nil {
		// If decryption fails, assume it's plain text
		return data, nil
	}
	return decrypted, nil
}

// InnovoLogin handles user login
// @Summary User login
// @Description Authenticate a user with email and password. Returns JWT token valid for 30 days.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body EncryptedLoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.APIResponse,jwt_token=string} "Login successful with JWT token"
// @Failure 400 {object} models.APIResponse
// @Router /Services/innovologin [post]
func (c *AuthController) InnovoLogin(ctx *gin.Context) {
	var encryptedReq EncryptedLoginRequest
	if err := ctx.ShouldBindJSON(&encryptedReq); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Initialize encryption service
	encryptService := services.NewEncryptDecryptService()

	// Safely decrypt or use plain text
	email, err := safeDecrypt(encryptService, encryptedReq.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to process email: " + err.Error(),
		})
		return
	}

	userpin, err := safeDecrypt(encryptService, encryptedReq.Userpin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to process userpin: " + err.Error(),
		})
		return
	}

	user, err := c.userService.CheckUser(email, userpin)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Invalid credentials",
			Response: 0,
		})
		return
	}

	// Check if account is deleted
	if user.UStatus == 5 {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Your account has been deleted. Please contact admin for more details.",
			Response: 0,
		})
		return
	}

	// Generate JWT token using Email and UserName
	jwtService := services.NewJWTService()
	token, err := jwtService.GenerateToken(user.Email, user.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to generate authentication token",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:        0,
		Message:     "OK",
		UserID:      &user.ID,
		UserDetails: user,
		JWTToken:    token,
	})
}

// EncryptedRegisterRequest represents the encrypted registration request
type EncryptedRegisterRequest struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	CNumber  *string `json:"cnumber"`
	Userpin  string  `json:"userpin"`
	Age      int     `json:"age"`
	Gender   string  `json:"gender"`
	Height   float64 `json:"height"`
	Weight   float64 `json:"weight"`
}

// InnovoRegister handles user registration
// @Summary User registration
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body EncryptedRegisterRequest true "Registration data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/innovoregister [post]
func (c *AuthController) InnovoRegister(ctx *gin.Context) {
	var encryptedReq EncryptedRegisterRequest
	if err := ctx.ShouldBindJSON(&encryptedReq); err != nil {
		fmt.Printf("DEBUG: JSON binding error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	fmt.Printf("DEBUG: Received encrypted request: %+v\n", encryptedReq)

	// Initialize encryption service
	encryptService := services.NewEncryptDecryptService()

	// Safely decrypt or use plain text for all fields
	email, err := safeDecrypt(encryptService, encryptedReq.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to process email: " + err.Error(),
		})
		return
	}

	username, err := safeDecrypt(encryptService, encryptedReq.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to process username: " + err.Error(),
		})
		return
	}

	userpin, err := safeDecrypt(encryptService, encryptedReq.Userpin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to process userpin: " + err.Error(),
		})
		return
	}

	// Process contact number if it exists
	var cNumber *string
	if encryptedReq.CNumber != nil {
		decrypted, err := safeDecrypt(encryptService, *encryptedReq.CNumber)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    1,
				Message: "Failed to process contact number: " + err.Error(),
			})
			return
		}
		cNumber = &decrypted
	}

	// Create the request
	req := models.RegisterRequest{
		Username: username,
		Email:    email,
		CNumber:  cNumber,
		Userpin:  userpin,
		Age:      encryptedReq.Age,
		Gender:   encryptedReq.Gender,
		Height:   encryptedReq.Height,
		Weight:   encryptedReq.Weight,
	}

	// Check if user already exists (service will handle encryption for lookup)
	existingUser, _ := c.userService.ValidateUser(req.Email)
	if existingUser != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User already exists with this email address",
			Response: 0,
		})
		return
	}

	// Register new user
	userID, err := c.userService.RegisterUser(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to register user",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "User registered successfully",
		UserID:  &userID,
	})
}
