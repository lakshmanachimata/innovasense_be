package services

import (
	"database/sql"
	"innovasense_be/config"
	"innovasense_be/models"
	"math"
	"time"
)

type HydrationService struct {
	db *sql.DB
}

func NewHydrationService() *HydrationService {
	return &HydrationService{
		db: config.GetDB(),
	}
}

// CalculateBMI calculates Body Mass Index
func (s *HydrationService) CalculateBMI(weight, height float64) float64 {
	heightInMeters := height / 100
	bmi := weight / (heightInMeters * heightInMeters)
	return math.Round(bmi*100) / 100
}

// CalculateTBSA calculates Total Body Surface Area using the Mosteller formula
func (s *HydrationService) CalculateTBSA(weight, height float64) float64 {
	tbsa := 0.007184 * math.Pow(height, 0.725) * math.Pow(weight, 0.425)
	return math.Round(tbsa*100) / 100
}

// CalculateSweatLoss calculates sweat loss based on device type
func (s *HydrationService) CalculateSweatLoss(tbsa, sweatPosition float64, deviceType int) float64 {
	var sweatLoss float64

	// Base calculation: ((58 * TBSA) - 73) * sweat_position
	baseCalculation := ((58 * tbsa) - 73) * sweatPosition

	// Apply device type multiplier
	if deviceType == 1 || deviceType == 3 {
		// Standard calculation
		sweatLoss = baseCalculation
	} else if deviceType == 2 || deviceType == 4 {
		// Double the calculation for enhanced devices
		sweatLoss = 2 * baseCalculation
	} else {
		// Default to standard calculation
		sweatLoss = baseCalculation
	}

	return math.Round(sweatLoss*100) / 100
}

// CalculateSweatRate calculates sweat rate in ml/hour
func (s *HydrationService) CalculateSweatRate(sweatLoss, timeTaken float64) float64 {
	// Convert to ml/hour: (sweat_loss / time_taken) * 60
	sweatRate := (sweatLoss / timeTaken) * 60
	return math.Round(sweatRate*100) / 100
}

// SaveHydrationData saves new hydration data with business logic calculations
func (s *HydrationService) SaveHydrationData(req *models.HydrationRequest) (int, error) {
	// Calculate BMI and TBSA
	bmi := s.CalculateBMI(req.Weight, req.Height)
	tbsa := s.CalculateTBSA(req.Weight, req.Height)

	// Calculate sweat loss based on device type
	sweatLoss := s.CalculateSweatLoss(tbsa, req.SweatPosition, req.DeviceType)

	// Calculate sweat rate
	sweatRate := s.CalculateSweatRate(sweatLoss, req.TimeTaken)

	query := `
		INSERT INTO user_data (user_id, weight, height, sweat_position, time_taken, 
		                      bmi, tbsa, image_path, sweat_rate, sweat_loss, device_type, 
		                      image_id, creation_datetime)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := s.db.Exec(query, req.UserID, req.Weight, req.Height, req.SweatPosition,
		req.TimeTaken, bmi, tbsa, req.ImagePath, sweatRate, sweatLoss,
		req.DeviceType, req.ImageID, now)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// SaveEnhancedHydrationData saves enhanced hydration data with additional calculations and summaries
func (s *HydrationService) SaveEnhancedHydrationData(req *models.HydrationRequest) (*models.EnhancedHydrationResponse, error) {
	// Calculate BMI and TBSA
	bmi := s.CalculateBMI(req.Weight, req.Height)
	tbsa := s.CalculateTBSA(req.Weight, req.Height)

	// Calculate sweat loss based on device type
	sweatLoss := s.CalculateSweatLoss(tbsa, req.SweatPosition, req.DeviceType)

	// Calculate sweat rate
	sweatRate := s.CalculateSweatRate(sweatLoss, req.TimeTaken)

	query := `
		INSERT INTO user_data (user_id, weight, height, sweat_position, time_taken, 
		                      bmi, tbsa, image_path, sweat_rate, sweat_loss, device_type, 
		                      image_id, creation_datetime)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := s.db.Exec(query, req.UserID, req.Weight, req.Height, req.SweatPosition,
		req.TimeTaken, bmi, tbsa, req.ImagePath, sweatRate, sweatLoss,
		req.DeviceType, req.ImageID, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get the saved data
	savedData, err := s.GetHydrationDataByID(int(id))
	if err != nil {
		return nil, err
	}

	// Get sweat summary if image_id is provided
	var sweatSummary []models.SweatImage
	if req.ImageID != nil {
		sweatSummary, _ = s.GetSweatSummaryByImageID(*req.ImageID)
	}

	// Get sweat rate summary
	sweatRateSummary, _ := s.GetSweatRateSummary(sweatRate)

	return &models.EnhancedHydrationResponse{
		ID:               int(id),
		Data:             savedData,
		SweatSummary:     sweatSummary,
		SweatRateSummary: sweatRateSummary,
	}, nil
}

// GetHydrationDataByID retrieves hydration data by ID
func (s *HydrationService) GetHydrationDataByID(id int) (*models.HydrationData, error) {
	query := `
		SELECT id, user_id, weight, height, sweat_position, time_taken, bmi, tbsa,
		       image_path, sweat_rate, sweat_loss, device_type, image_id, creation_datetime
		FROM user_data 
		WHERE id = ?
	`

	var data models.HydrationData
	var creationDatetimeStr string
	err := s.db.QueryRow(query, id).Scan(&data.ID, &data.UserID, &data.Weight, &data.Height,
		&data.SweatPosition, &data.TimeTaken, &data.BMI, &data.TBSA, &data.ImagePath,
		&data.SweatRate, &data.SweatLoss, &data.DeviceType, &data.ImageID, &creationDatetimeStr)
	if err != nil {
		return nil, err
	}

	// Parse creation_datetime from string to time.Time
	if creationDatetimeStr != "" {
		// Try multiple time formats
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02 15:04:05.000",
		}

		for _, format := range formats {
			if parsedTime, err := time.Parse(format, creationDatetimeStr); err == nil {
				data.CreationDatetime = parsedTime
				break
			}
		}
	}

	return &data, nil
}

// GetSweatSummaryByImageID retrieves sweat summary by image ID
func (s *HydrationService) GetSweatSummaryByImageID(imageID int) ([]models.SweatImage, error) {
	query := `
		SELECT id, image_path, sweat_range, implications, recomm, strategy, result, colorcode
		FROM sweat_images 
		WHERE id = ?
	`

	rows, err := s.db.Query(query, imageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.SweatImage
	for rows.Next() {
		var image models.SweatImage
		err := rows.Scan(&image.ID, &image.ImagePath, &image.SweatRange, &image.Implications,
			&image.Recomm, &image.Strategy, &image.Result, &image.ColorCode)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

// GetSweatRateSummary retrieves sweat rate summary based on calculated sweat rate
func (s *HydrationService) GetSweatRateSummary(sweatRate float64) ([]models.SweatRateSummary, error) {
	query := `
		SELECT id, low_limit, high_limit, hyd_status, comments, recomm, color
		FROM sweatrate_summary 
		WHERE ? BETWEEN low_limit AND high_limit
	`

	rows, err := s.db.Query(query, sweatRate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.SweatRateSummary
	for rows.Next() {
		var item models.SweatRateSummary
		err := rows.Scan(&item.ID, &item.LowLimit, &item.HighLimit, &item.HydStatus,
			&item.Comments, &item.Recomm, &item.Color)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

// UpdateHydrationData updates existing hydration data
func (s *HydrationService) UpdateHydrationData(req *models.UpdateHydrationRequest) error {
	query := `
		UPDATE user_data 
		SET weight = ?, height = ?, sweat_position = ?, time_taken = ?, 
		    bmi = ?, tbsa = ?, sweat_rate = ?, sweat_loss = ?, device_type = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(query, req.Weight, req.Height, req.SweatPosition, req.TimeTaken,
		req.BMI, req.TBSA, req.SweatRate, req.SweatLoss, req.DeviceType, req.ID)
	return err
}

// UpdateSweatData updates sweat analysis data
func (s *HydrationService) UpdateSweatData(req *models.UpdateSweatDataRequest) error {
	query := `
		UPDATE user_data 
		SET sweat_rate = ?, sweat_loss = ?
		WHERE user_id = ? AND image_id = ?
	`

	_, err := s.db.Exec(query, req.SweatRate, req.SweatLoss, req.UserID, req.ImageID)
	return err
}

// GetHydrationHistory retrieves hydration history for a user (matches PHP logic)
func (s *HydrationService) GetHydrationHistory(userID int, fromDate, toDate string) ([]models.HydrationData, error) {
	// PHP adds +1 to to_date for inclusive range
	query := `
		SELECT id, user_id, weight, height, sweat_position, time_taken, bmi, tbsa,
		       image_path, sweat_rate, sweat_loss, device_type, image_id, creation_datetime
		FROM user_data 
		WHERE user_id = ? AND creation_datetime >= ? AND creation_datetime <= ?
		ORDER BY creation_datetime DESC
	`

	// Add +1 day to to_date for inclusive range (matching PHP logic)
	rows, err := s.db.Query(query, userID, fromDate, toDate+" 23:59:59")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.HydrationData
	for rows.Next() {
		var item models.HydrationData
		var creationDatetimeStr string
		err := rows.Scan(&item.ID, &item.UserID, &item.Weight, &item.Height, &item.SweatPosition,
			&item.TimeTaken, &item.BMI, &item.TBSA, &item.ImagePath, &item.SweatRate,
			&item.SweatLoss, &item.DeviceType, &item.ImageID, &creationDatetimeStr)
		if err != nil {
			return nil, err
		}

		// Parse creation_datetime from string to time.Time
		if creationDatetimeStr != "" {
			// Try multiple time formats
			formats := []string{
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05.000Z",
				"2006-01-02 15:04:05.000",
			}

			for _, format := range formats {
				if parsedTime, err := time.Parse(format, creationDatetimeStr); err == nil {
					item.CreationDatetime = parsedTime
					break
				}
			}
		}

		data = append(data, item)
	}

	return data, nil
}

// GetElectrolyteHistory retrieves electrolyte history (matches PHP logic - only creation_datetime and image_id)
func (s *HydrationService) GetElectrolyteHistory(userID int, fromDate, toDate string) ([]models.ElectrolyteHistoryData, error) {
	// PHP only returns creation_datetime and image_id for electrolyte history
	query := `
		SELECT creation_datetime, image_id
		FROM user_data 
		WHERE user_id = ? AND device_type IN (3, 4) AND creation_datetime >= ? AND creation_datetime <= ?
		ORDER BY creation_datetime DESC
	`

	// Add +1 day to to_date for inclusive range (matching PHP logic)
	rows, err := s.db.Query(query, userID, fromDate, toDate+" 23:59:59")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.ElectrolyteHistoryData
	for rows.Next() {
		var item models.ElectrolyteHistoryData
		var creationDatetimeStr string
		err := rows.Scan(&creationDatetimeStr, &item.ImageID)
		if err != nil {
			return nil, err
		}

		// Parse creation_datetime from string to time.Time
		if creationDatetimeStr != "" {
			// Try multiple time formats
			formats := []string{
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05.000Z",
				"2006-01-02 15:04:05.000",
			}

			for _, format := range formats {
				if parsedTime, err := time.Parse(format, creationDatetimeStr); err == nil {
					item.CreationDatetime = parsedTime
					break
				}
			}
		}

		data = append(data, item)
	}

	return data, nil
}

// GetSummary retrieves summary data based on sweat position (matches PHP logic - uses sweat_summary table)
func (s *HydrationService) GetSummary(sweatPosition float64) ([]models.SweatSummary, error) {
	// PHP uses sweat_summary table, not sweatrate_summary
	query := `
		SELECT id, low_limit, high_limit, hyd_status, comments, recomm, color
		FROM sweat_summary 
		WHERE ? BETWEEN low_limit AND high_limit
	`

	rows, err := s.db.Query(query, sweatPosition)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.SweatSummary
	for rows.Next() {
		var item models.SweatSummary
		err := rows.Scan(&item.ID, &item.LowLimit, &item.HighLimit, &item.HydStatus,
			&item.Comments, &item.Recomm, &item.Color)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

// GetUserDetailedSummary retrieves detailed summary for a user (matches PHP logic with multiple calls)
func (s *HydrationService) GetUserDetailedSummary(id int) (*models.DetailedSummaryResponse, error) {
	// Get the main hydration data
	hydrationData, err := s.GetHydrationDataByID(id)
	if err != nil {
		return nil, err
	}

	// Get summary based on sweat position
	summary, _ := s.GetSummary(hydrationData.SweatPosition)

	// Get sweat summary if image_id is provided
	var sweatSummary []models.SweatImage
	if hydrationData.ImageID != nil {
		sweatSummary, _ = s.GetSweatSummaryByImageID(*hydrationData.ImageID)
	}

	// Get sweat rate summary
	sweatRateSummary, _ := s.GetSweatRateSummary(hydrationData.SweatRate)

	return &models.DetailedSummaryResponse{
		Data:             hydrationData,
		Summary:          summary,
		SweatSummary:     sweatSummary,
		SweatRateSummary: sweatRateSummary,
	}, nil
}

// GetHydrationSummaryScreen retrieves formatted data for the summary screen
func (s *HydrationService) GetHydrationSummaryScreen(id int) (*models.HydrationSummaryData, error) {
	// Get the main hydration data
	hydrationData, err := s.GetHydrationDataByID(id)
	if err != nil {
		return nil, err
	}

	// Get summary based on sweat position (not used in current implementation)
	_, _ = s.GetSummary(hydrationData.SweatPosition)

	// Get sweat summary if image_id is provided
	var sweatSummary []models.SweatSummaryItem
	if hydrationData.ImageID != nil {
		sweatImages, _ := s.GetSweatSummaryByImageID(*hydrationData.ImageID)
		// Convert SweatImage to SweatSummaryItem
		for _, img := range sweatImages {
			sweatSummary = append(sweatSummary, models.SweatSummaryItem{
				ID:           img.ID,
				ImagePath:    img.ImagePath,
				SweatRange:   img.SweatRange,
				Implications: img.Implications,
				Recommend:    img.Recomm,
				Strategy:     img.Strategy,
				Result:       img.Result,
				ColorCode:    img.ColorCode,
			})
		}
	}

	// Get sweat rate summary
	sweatRateSummary, _ := s.GetSweatRateSummary(hydrationData.SweatRate)
	var sweatRateSummaryItems []models.SweatRateSummaryItem
	for _, rate := range sweatRateSummary {
		sweatRateSummaryItems = append(sweatRateSummaryItems, models.SweatRateSummaryItem{
			ID:        rate.ID,
			LowLimit:  rate.LowLimit,
			HighLimit: rate.HighLimit,
			HydStatus: rate.HydStatus,
			Comments:  rate.Comments,
			Recommend: rate.Recomm,
			Color:     rate.Color,
		})
	}

	// Create the summary data
	summaryData := &models.HydrationSummaryData{
		ID:               hydrationData.ID,
		UserID:           hydrationData.UserID,
		Weight:           hydrationData.Weight,
		Height:           hydrationData.Height,
		SweatPosition:    hydrationData.SweatPosition,
		TimeTaken:        hydrationData.TimeTaken,
		BMI:              hydrationData.BMI,
		TBSA:             hydrationData.TBSA,
		ImagePath:        *hydrationData.ImagePath,
		SweatRate:        hydrationData.SweatRate,
		SweatLoss:        hydrationData.SweatLoss,
		DeviceType:       hydrationData.DeviceType,
		ImageID:          *hydrationData.ImageID,
		CreationDatetime: hydrationData.CreationDatetime,
		SweatSummary:     sweatSummary,
		SweatRateSummary: sweatRateSummaryItems,
	}

	// Calculate additional display fields
	if len(sweatRateSummaryItems) > 0 {
		summaryData.HydrationStatus = sweatRateSummaryItems[0].HydStatus
		summaryData.RiskLevel = s.calculateRiskLevel(hydrationData.SweatRate)
		summaryData.Recommendations = sweatRateSummaryItems[0].Recommend
		summaryData.NextTestDate = s.calculateNextTestDate(hydrationData.CreationDatetime)
	}

	return summaryData, nil
}

// calculateRiskLevel determines the risk level based on sweat rate
func (s *HydrationService) calculateRiskLevel(sweatRate float64) string {
	if sweatRate < 200 {
		return "Low"
	} else if sweatRate < 500 {
		return "Moderate"
	} else {
		return "High"
	}
}

// calculateNextTestDate suggests when the next test should be done
func (s *HydrationService) calculateNextTestDate(lastTestDate time.Time) string {
	// Suggest next test in 7 days for high risk, 14 days for moderate, 30 days for low
	nextDate := lastTestDate.AddDate(0, 0, 7) // Default to 7 days
	return nextDate.Format("2006-01-02")
}

// GetClientHistory retrieves client history (matches PHP logic)
func (s *HydrationService) GetClientHistory(userID int) ([]models.HydrationData, error) {
	query := `
		SELECT id, user_id, weight, height, sweat_position, time_taken, bmi, tbsa,
		       image_path, sweat_rate, sweat_loss, device_type, image_id, creation_datetime
		FROM user_data 
		WHERE user_id = ?
		ORDER BY creation_datetime DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.HydrationData
	for rows.Next() {
		var item models.HydrationData
		var creationDatetimeStr string
		err := rows.Scan(&item.ID, &item.UserID, &item.Weight, &item.Height, &item.SweatPosition,
			&item.TimeTaken, &item.BMI, &item.TBSA, &item.ImagePath, &item.SweatRate,
			&item.SweatLoss, &item.DeviceType, &item.ImageID, &creationDatetimeStr)
		if err != nil {
			return nil, err
		}

		// Parse creation_datetime from string to time.Time
		if creationDatetimeStr != "" {
			// Try multiple time formats
			formats := []string{
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05.000Z",
				"2006-01-02 15:04:05.000",
			}

			for _, format := range formats {
				if parsedTime, err := time.Parse(format, creationDatetimeStr); err == nil {
					item.CreationDatetime = parsedTime
					break
				}
			}
		}

		data = append(data, item)
	}

	return data, nil
}
