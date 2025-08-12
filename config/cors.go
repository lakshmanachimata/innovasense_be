package config

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// GetCORSConfig returns CORS configuration based on environment
func GetCORSConfig() *CORSConfig {
	env := os.Getenv("GIN_MODE")
	
	if env == "release" {
		// Production CORS - more restrictive
		return &CORSConfig{
			AllowedOrigins:   []string{"https://yourdomain.com", "https://www.yourdomain.com"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
			ExposedHeaders:   []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           86400, // 24 hours
		}
	}
	
	// Development/Test CORS - more permissive
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With", "Content-Range", "Content-Disposition", "Content-Description"},
		ExposedHeaders:   []string{"Content-Length", "Content-Range", "Content-Disposition", "Content-Description"},
		AllowCredentials: true,
		MaxAge:           3600, // 1 hour
	}
}

// CORSMiddleware returns a configured CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	config := GetCORSConfig()
	
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		if len(config.AllowedOrigins) > 0 && config.AllowedOrigins[0] != "*" {
			allowed := false
			for _, allowedOrigin := range config.AllowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}
			if !allowed {
				c.Header("Access-Control-Allow-Origin", "")
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
		c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
