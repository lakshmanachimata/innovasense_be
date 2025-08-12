package controllers

import (
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

// InnovoLogin handles user login
// @Summary User login
// @Description Authenticate a user with contact number and password. Returns JWT token valid for 30 days.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.APIResponse,jwt_token=string} "Login successful with JWT token"
// @Failure 400 {object} models.APIResponse
// @Router /Services/innovologin [post]
func (c *AuthController) InnovoLogin(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	user, err := c.userService.CheckUser(req.CNumber, req.Userpin)
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

	// Generate JWT token using CNumber and UserName
	jwtService := services.NewJWTService()
	token, err := jwtService.GenerateToken(user.CNumber, user.Username)
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

// InnovoRegister handles user registration
// @Summary User registration
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/innovoregister [post]
func (c *AuthController) InnovoRegister(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Check if user already exists
	existingUser, _ := c.userService.ValidateUser(req.CNumber)
	if existingUser != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User already exists with this contact number",
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
