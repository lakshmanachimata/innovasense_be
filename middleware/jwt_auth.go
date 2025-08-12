package middleware

import (
	"innovasense_be/models"
	"innovasense_be/services"
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
		c.Set("user_cnumber", claims.CNumber)
		c.Set("username", claims.UserName)

		c.Next()
	}
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

// GetUserCNumberFromJWTContext retrieves user CNumber from JWT context
func GetUserCNumberFromJWTContext(c *gin.Context) (string, bool) {
	cNumber, exists := c.Get("user_cnumber")
	if !exists {
		return "", false
	}

	if cn, ok := cNumber.(string); ok {
		return cn, true
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
