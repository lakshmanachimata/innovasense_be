package models

import "time"

// HydrationSummaryViewModel represents the view model for the hydration summary screen
type HydrationSummaryViewModel struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    *HydrationSummaryData `json:"data"`
}

// HydrationSummaryData represents the main data for the summary screen
type HydrationSummaryData struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	Weight           float64   `json:"weight"`
	Height           float64   `json:"height"`
	SweatPosition    float64   `json:"sweat_position"`
	TimeTaken        float64   `json:"time_taken"`
	BMI              float64   `json:"bmi"`
	TBSA             float64   `json:"tbsa"`
	ImagePath        string    `json:"image_path"`
	SweatRate        float64   `json:"sweat_rate"`
	SweatLoss        float64   `json:"sweat_loss"`
	DeviceType       int       `json:"device_type"`
	ImageID          int       `json:"image_id"`
	CreationDatetime time.Time `json:"creation_datetime"`

	// Summary sections
	SweatSummary     []SweatSummaryItem     `json:"sweat_summary"`
	SweatRateSummary []SweatRateSummaryItem `json:"sweat_rate_summary"`

	// Additional calculated fields for display
	HydrationStatus string `json:"hydration_status"`
	RiskLevel       string `json:"risk_level"`
	Recommendations string `json:"recommendations"`
	NextTestDate    string `json:"next_test_date"`
}

// SweatSummaryItem represents individual sweat summary items
type SweatSummaryItem struct {
	ID           int    `json:"id"`
	ImagePath    string `json:"image_path"`
	SweatRange   string `json:"sweat_range"`
	Implications string `json:"implications"`
	Recommend    string `json:"recomm"`
	Strategy     string `json:"strategy"`
	Result       string `json:"result"`
	ColorCode    string `json:"colorcode"`
}

// SweatRateSummaryItem represents individual sweat rate summary items
type SweatRateSummaryItem struct {
	ID        int     `json:"id"`
	LowLimit  float64 `json:"low_limit"`
	HighLimit float64 `json:"high_limit"`
	HydStatus string  `json:"hyd_status"`
	Comments  string  `json:"comments"`
	Recommend string  `json:"recomm"`
	Color     string  `json:"color"`
}

// HydrationSummaryRequest represents the request for getting hydration summary
type HydrationSummaryRequest struct {
	CNumber  string `json:"cnumber" binding:"required"`
	Username string `json:"username" binding:"required"`
	ID       int    `json:"id" binding:"required"`
}

// HydrationSummaryResponse represents the complete response for the summary screen
type HydrationSummaryResponse struct {
	Code     int                   `json:"code"`
	Message  string                `json:"message"`
	Response *HydrationSummaryData `json:"response"`
}
