package models

// BannerImage represents banner images
type BannerImage struct {
	ID        int    `json:"id" db:"id"`
	ImagePath string `json:"image_path" db:"image_path"`
}

// HomeImage represents home images
type HomeImage struct {
	ID        int    `json:"id" db:"id"`
	ImagePath string `json:"image_path" db:"image_path"`
}

// SweatImage represents sweat analysis images
type SweatImage struct {
	ID           int    `json:"id" db:"id"`
	ImagePath    string `json:"image_path" db:"image_path"`
	SweatRange   string `json:"sweat_range" db:"sweat_range"`
	Implications string `json:"implications" db:"implications"`
	Recomm       string `json:"recomm" db:"recomm"`
	Strategy     string `json:"strategy" db:"strategy"`
	Result       string `json:"result" db:"result"`
	ColorCode    string `json:"colorcode" db:"colorcode"`
}

// Device represents device master data
type Device struct {
	ID         int    `json:"id" db:"id"`
	DeviceName string `json:"device_name" db:"device_name"`
	DeviceText string `json:"device_text" db:"device_text"`
}

// SweatRateSummary represents sweat rate summary data
type SweatRateSummary struct {
	ID        int     `json:"id" db:"id"`
	LowLimit  float64 `json:"low_limit" db:"low_limit"`
	HighLimit float64 `json:"high_limit" db:"high_limit"`
	HydStatus string  `json:"hyd_status" db:"hyd_status"`
	Comments  string  `json:"comments" db:"comments"`
	Recomm    string  `json:"recomm" db:"recomm"`
	Color     string  `json:"color" db:"color"`
}

// ImageUploadRequest represents image upload request
type ImageUploadRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	// Note: ImagePath will be generated on the server after file upload
}

// UpdateImagePathRequest represents image path update request
type UpdateImagePathRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required"`
	UserID    int    `json:"userid" binding:"required"`
	ImageID   int    `json:"image_id" binding:"required"`
	ImagePath string `json:"image_path" binding:"required"`
}

// GetSweatImagesRequest represents sweat images retrieval request
type GetSweatImagesRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
}
