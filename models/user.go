package models

import (
	"time"
)

// User represents the users_master table
// @Description User data structure from users_master table
type User struct {
	// @Description Unique user identifier
	ID int `json:"id" db:"id"`
	// @Description User's email address
	Email string `json:"email" db:"email"`
	// @Description User's contact number (optional)
	CNumber *string `json:"cnumber,omitempty" db:"cnumber"`
	// @Description User's encrypted PIN/password
	Userpin string `json:"userpin" db:"userpin"`
	// @Description User's display name
	Username string `json:"username" db:"username"`
	// @Description User's gender
	Gender string `json:"gender" db:"gender"`
	// @Description User's age in years
	Age int `json:"age" db:"age"`
	// @Description User's height in cm
	Height float64 `json:"height" db:"height"`
	// @Description User's weight in kg
	Weight float64 `json:"weight" db:"weight"`
	// @Description User's role identifier
	RoleID int `json:"role_id" db:"role_id"`
	// @Description User's status (0=active, 1=inactive)
	UStatus int `json:"ustatus" db:"ustatus"`
	// @Description User account creation timestamp
	CreationDatetime time.Time `json:"creation_datetime" db:"creation_datetime"`
}

// LoginRequest represents login request data
// @Description Login request data structure
type LoginRequest struct {
	// @Description User's email address (required)
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	// @Description User's PIN/password (required)
	Userpin string `json:"userpin" binding:"required" example:"test@123"`
}

// RegisterRequest represents user registration request
// @Description User registration request data structure
type RegisterRequest struct {
	// @Description User's email address (required)
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	// @Description User's contact number (optional)
	CNumber *string `json:"cnumber,omitempty" example:"1234567890"`
	// @Description User's PIN/password (required)
	Userpin string `json:"userpin" binding:"required" example:"test@123"`
	// @Description User's display name (required)
	Username string `json:"username" binding:"required" example:"John Doe"`
	// @Description User's gender
	Gender string `json:"gender" example:"Male"`
	// @Description User's age in years
	Age int `json:"age" example:"25"`
	// @Description User's height in cm
	Height float64 `json:"height" example:"170.5"`
	// @Description User's weight in kg
	Weight float64 `json:"weight" example:"70.0"`
}

// APIResponse represents the standard API response format
// @Description Standard API response format used across all endpoints
type APIResponse struct {
	// @Description Response code: 0 for success, 1 for error
	Code int `json:"code" example:"0"`
	// @Description Response message
	Message string `json:"message" example:"OK"`
	// @Description Main response data (varies by endpoint)
	Response interface{} `json:"response,omitempty"`
	// @Description User ID when applicable
	UserID *int `json:"userid,omitempty" example:"123"`
	// @Description User details when applicable
	UserDetails interface{} `json:"userdetails,omitempty"`

	// @Description Users list when applicable
	Users interface{} `json:"users,omitempty"`
	// @Description Summary data when applicable
	Summary interface{} `json:"summary,omitempty"`
	// @Description Sweat summary when applicable
	SweatSummary interface{} `json:"sweatsummary,omitempty"`
	// @Description Sweat rate summary when applicable
	SweatRateSummary interface{} `json:"sweatratesummary,omitempty"`
	// @Description JWT token when applicable
	JWTToken string `json:"jwt_token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
