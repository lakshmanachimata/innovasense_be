package controllers

import (
	"fmt"
	"innovasense_be/models"
	"innovasense_be/services"
	"net/http"

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

	// Decrypt the encrypted fields
	decryptedEmail, err := encryptService.GetDecryptData(encryptedReq.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email: " + err.Error(),
		})
		return
	}

	decryptedUserpin, err := encryptService.GetDecryptData(encryptedReq.Userpin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt userpin: " + err.Error(),
		})
		return
	}

	user, err := c.userService.CheckUser(decryptedEmail, decryptedUserpin)
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

	// Decrypt the encrypted fields
	decryptedEmail, err := encryptService.GetDecryptData(encryptedReq.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt email: " + err.Error(),
		})
		return
	}

	decryptedUsername, err := encryptService.GetDecryptData(encryptedReq.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt username: " + err.Error(),
		})
		return
	}

	decryptedUserpin, err := encryptService.GetDecryptData(encryptedReq.Userpin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to decrypt userpin: " + err.Error(),
		})
		return
	}

	// Decrypt contact number if it exists
	var decryptedCNumber *string
	if encryptedReq.CNumber != nil {
		decrypted, err := encryptService.GetDecryptData(*encryptedReq.CNumber)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    1,
				Message: "Failed to decrypt contact number: " + err.Error(),
			})
			return
		}
		decryptedCNumber = &decrypted
	}

	// Create the decrypted request
	req := models.RegisterRequest{
		Username: decryptedUsername,
		Email:    decryptedEmail,
		CNumber:  decryptedCNumber,
		Userpin:  decryptedUserpin,
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
