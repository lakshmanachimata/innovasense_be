package middleware

import (
	"bytes"
	"encoding/json"
	"innovasense_be/models"
	"innovasense_be/services"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates JWT tokens for protected endpoints
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Code:    1,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Code:    1,
				Message: "Invalid authorization format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Code:    1,
				Message: "JWT token is required",
			})
			c.Abort()
			return
		}

		// Validate the JWT token
		jwtService := services.NewJWTService()
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Code:    1,
				Message: "Invalid or expired JWT token",
			})
			c.Abort()
			return
		}

		// Set claims in context for later use
		c.Set("jwt_claims", claims)
		c.Set("user_email", claims.Email)
		c.Set("username", claims.UserName)

		// Validate request body cnumber and username against JWT claims
		if err := validateUserIdentity(c, claims); err != nil {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Code:    1,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// validateUserIdentity checks if the request body contains email and username that match JWT claims
func validateUserIdentity(c *gin.Context, claims *services.Claims) error {
	// Only validate for POST requests with JSON body
	if c.Request.Method != "POST" {
		return nil
	}

	// Check if request has JSON content
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	// Restore request body for later use
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Try to parse as common request structure
	var commonRequest struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	if err := json.Unmarshal(body, &commonRequest); err != nil {
		// If parsing fails, it might be a different request structure
		// We'll let the controller handle validation
		return nil
	}

	// Validate email and username against JWT claims
	if commonRequest.Email != "" && commonRequest.Email != claims.Email {
		return &ValidationError{Message: "email in request body does not match authenticated user"}
	}

	if commonRequest.Username != "" && commonRequest.Username != claims.UserName {
		return &ValidationError{Message: "username in request body does not match authenticated user"}
	}

	return nil
}

// ValidationError represents validation errors
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// GetJWTClaimsFromContext retrieves JWT claims from gin context
func GetJWTClaimsFromContext(c *gin.Context) (*services.Claims, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}

	if jwtClaims, ok := claims.(*services.Claims); ok {
		return jwtClaims, true
	}

	return nil, false
}

// GetUserEmailFromJWTContext retrieves user Email from JWT context
func GetUserEmailFromJWTContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	if em, ok := email.(string); ok {
		return em, true
	}

	return "", false
}

// GetUserNameFromJWTContext retrieves username from JWT context
func GetUserNameFromJWTContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	if un, ok := username.(string); ok {
		return un, true
	}

	return "", false
}
