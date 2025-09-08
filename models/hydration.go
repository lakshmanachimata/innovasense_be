package models

import (
	"time"
)

// HydrationData represents the user_data table for hydration records
type HydrationData struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	Weight           float64   `json:"weight" db:"weight"`
	Height           float64   `json:"height" db:"height"`
	SweatPosition    float64   `json:"sweat_position" db:"sweat_position"`
	TimeTaken        float64   `json:"time_taken" db:"time_taken"`
	BMI              float64   `json:"bmi" db:"bmi"`
	TBSA             float64   `json:"tbsa" db:"tbsa"`
	ImagePath        *string   `json:"image_path" db:"image_path"`
	SweatRate        float64   `json:"sweat_rate" db:"sweat_rate"`
	SweatLoss        float64   `json:"sweat_loss" db:"sweat_loss"`
	DeviceType       int       `json:"device_type" db:"device_type"`
	ImageID          *int      `json:"image_id" db:"image_id"`
	CreationDatetime time.Time `json:"creation_datetime" db:"creation_datetime"`
}

// EnhancedHydrationResponse represents the enhanced hydration response with summaries
type EnhancedHydrationResponse struct {
	ID               int                `json:"id"`
	Data             *HydrationData     `json:"data"`
	SweatSummary     []SweatImage       `json:"sweatsummary"`
	SweatRateSummary []SweatRateSummary `json:"sweatratesummary"`
}

// DetailedSummaryResponse represents the detailed summary response with multiple summaries
type DetailedSummaryResponse struct {
	Data             *HydrationData     `json:"data"`
	Summary          []SweatSummary     `json:"summary"`
	SweatSummary     []SweatImage       `json:"SweatSummary"`
	SweatRateSummary []SweatRateSummary `json:"SweatRateSummary"`
}

// ElectrolyteHistoryData represents electrolyte history data (matches PHP - only creation_datetime and image_id)
type ElectrolyteHistoryData struct {
	CreationDatetime time.Time `json:"creation_datetime" db:"creation_datetime"`
	ImageID          *int      `json:"image_id" db:"image_id"`
}

// SweatSummary represents sweat summary data from sweat_summary table (matches PHP logic)
type SweatSummary struct {
	ID        int     `json:"id" db:"id"`
	LowLimit  float64 `json:"low_limit" db:"low_limit"`
	HighLimit float64 `json:"high_limit" db:"high_limit"`
	HydStatus string  `json:"hyd_status" db:"hyd_status"`
	Comments  string  `json:"comments" db:"comments"`
	Recomm    string  `json:"recomm" db:"recomm"`
	Color     string  `json:"color" db:"color"`
}

// HydrationRequest represents hydration data submission request
type HydrationRequest struct {
	Email         string  `json:"email" binding:"required,email"`
	Username      string  `json:"username" binding:"required"`
	UserID        int     `json:"userid" binding:"required"`
	Weight        float64 `json:"weight" binding:"required"`
	Height        float64 `json:"height" binding:"required"`
	SweatPosition float64 `json:"sweat_position" binding:"required"`
	TimeTaken     float64 `json:"time_taken" binding:"required"`
	ImagePath     *string `json:"image_path"`
	DeviceType    int     `json:"device_type" binding:"required"`
	ImageID       *int    `json:"image_id"`
}

// UpdateHydrationRequest represents hydration data update request
type UpdateHydrationRequest struct {
	Email         string  `json:"email" binding:"required,email"`
	Username      string  `json:"username" binding:"required"`
	ID            int     `json:"id" binding:"required"`
	Weight        float64 `json:"weight"`
	Height        float64 `json:"height"`
	SweatPosition float64 `json:"sweat_position"`
	TimeTaken     float64 `json:"time_taken"`
	BMI           float64 `json:"bmi"`
	TBSA          float64 `json:"tbsa"`
	SweatRate     float64 `json:"sweat_rate"`
	SweatLoss     float64 `json:"sweat_loss"`
	DeviceType    int     `json:"device_type"`
}

// SweatData represents sweat analysis data
type SweatData struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	ImageID          int       `json:"image_id" db:"image_id"`
	SweatRate        float64   `json:"sweat_rate" db:"sweat_rate"`
	SweatLoss        float64   `json:"sweat_loss" db:"sweat_loss"`
	CreationDatetime time.Time `json:"creation_datetime" db:"creation_datetime"`
}

// UpdateSweatDataRequest represents sweat data update request
type UpdateSweatDataRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Username  string  `json:"username" binding:"required"`
	UserID    int     `json:"userid" binding:"required"`
	ImageID   int     `json:"image_id" binding:"required"`
	SweatRate float64 `json:"sweat_rate" binding:"required"`
	SweatLoss float64 `json:"sweat_loss" binding:"required"`
}

// HistoryRequest represents history data request
type HistoryRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	UserID   int    `json:"userid" binding:"required"`
	FromDate string `json:"from_date" binding:"required"`
	ToDate   string `json:"to_date" binding:"required"`
}

// SummaryRequest represents summary data request
type SummaryRequest struct {
	Email         string  `json:"email" binding:"required,email"`
	Username      string  `json:"username" binding:"required"`
	SweatPosition float64 `json:"sweat_position" binding:"required"`
}

// DetailedSummaryRequest represents detailed summary request
type DetailedSummaryRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	ID       int    `json:"id" binding:"required"`
}

// ClientHistoryRequest represents client history request
type ClientHistoryRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
}

// HydrationRecommendationRequest represents the request for hydration recommendation
type HydrationRecommendationRequest struct {
	Name          string  `json:"name" binding:"required"`
	Contact       string  `json:"contact" binding:"required"` // Contact number or email
	Gender        string  `json:"gender" binding:"required"`
	Age           int     `json:"age" binding:"required"`
	SweatPosition float64 `json:"sweat_position" binding:"required"`
	WorkoutTime   float64 `json:"workout_time" binding:"required"`
	Height        float64 `json:"height" binding:"required"`
	Weight        float64 `json:"weight" binding:"required"`
}

// OrgAuthRequest represents organization authentication request
type OrgAuthRequest struct {
	APIKey    string `json:"apikey" binding:"required"`
	SecretKey string `json:"secretkey" binding:"required"`
}

// Organization represents organization data
type Organization struct {
	ID      int    `json:"id" db:"id"`
	OrgName string `json:"org_name" db:"org_name"`
	OrgDesc string `json:"org_desc" db:"org_desc"`
	SaltKey string `json:"salt_key" db:"salt_key"`
	APIKey  string `json:"api_key" db:"api_key"`
}

// OrgUser represents organization user data
type OrgUser struct {
	ID       int    `json:"id" db:"id"`
	EmailID  string `json:"email_id" db:"email_id"`
	UserPwd  string `json:"user_pwd" db:"user_pwd"`
	UserName string `json:"user_name" db:"user_name"`
	OrgID    int    `json:"org_id" db:"org_id"`
}

// HistoricalDataRequest represents the request for historical data
type HistoricalDataRequest struct {
	Contact  string `json:"contact" binding:"required"` // Contact number or email
	FromDate string `json:"from_date"`                  // Optional from date
	ToDate   string `json:"to_date"`                    // Optional to date
}

// HistoricalDataResponse represents the historical data response
type HistoricalDataResponse struct {
	SweatPosition []HistoricalDataItem `json:"sweat_position"`
	SweatRate     []HistoricalDataItem `json:"sweat_rate"`
	SweatLoss     []HistoricalDataItem `json:"sweat_loss"`
}

// HistoricalDataItem represents a single historical data item
type HistoricalDataItem struct {
	Datetime string  `json:"datetime"`
	Value    float64 `json:"value"`
}
