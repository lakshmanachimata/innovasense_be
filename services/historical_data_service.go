package services

import (
	"database/sql"
	"fmt"
	"innovasense_be/models"
	"time"
)

type HistoricalDataService struct {
	db *sql.DB
}

func NewHistoricalDataService(db *sql.DB) *HistoricalDataService {
	return &HistoricalDataService{db: db}
}

// GetHistoricalData retrieves historical data for a user
func (s *HistoricalDataService) GetHistoricalData(userID int, fromDate, toDate string) (*models.HistoricalDataResponse, error) {
	var whereClause string
	var args []interface{}

	if fromDate != "" && toDate != "" {
		whereClause = "WHERE user_id = ? AND DATE(creation_datetime) BETWEEN ? AND ?"
		args = []interface{}{userID, fromDate, toDate}
	} else {
		whereClause = "WHERE user_id = ?"
		args = []interface{}{userID}
	}

	// Get sweat position data
	sweatPositionQuery := fmt.Sprintf(`
		SELECT creation_datetime, sweat_position 
		FROM user_data 
		%s 
		ORDER BY creation_datetime DESC
	`, whereClause)

	sweatPositionRows, err := s.db.Query(sweatPositionQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sweat position data: %v", err)
	}
	defer sweatPositionRows.Close()

	var sweatPosition []models.HistoricalDataItem
	for sweatPositionRows.Next() {
		var item models.HistoricalDataItem
		var creationDatetimeStr string
		err := sweatPositionRows.Scan(&creationDatetimeStr, &item.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sweat position data: %v", err)
		}

		// Parse the timestamp string with multiple format attempts
		var creationDatetime time.Time
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000000Z",
			"2006-01-02 15:04:05.000000",
		}

		parsed := false
		for _, format := range formats {
			if t, err := time.Parse(format, creationDatetimeStr); err == nil {
				creationDatetime = t
				parsed = true
				break
			}
		}

		if !parsed {
			// If parsing fails, use the original string
			item.Datetime = creationDatetimeStr
		} else {
			item.Datetime = creationDatetime.Format("2006-01-02 15:04:05")
		}

		sweatPosition = append(sweatPosition, item)
	}

	// Get sweat rate data
	sweatRateQuery := fmt.Sprintf(`
		SELECT creation_datetime, sweat_rate 
		FROM user_data 
		%s 
		ORDER BY creation_datetime DESC
	`, whereClause)

	sweatRateRows, err := s.db.Query(sweatRateQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sweat rate data: %v", err)
	}
	defer sweatRateRows.Close()

	var sweatRate []models.HistoricalDataItem
	for sweatRateRows.Next() {
		var item models.HistoricalDataItem
		var creationDatetimeStr string
		err := sweatRateRows.Scan(&creationDatetimeStr, &item.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sweat rate data: %v", err)
		}

		// Parse the timestamp string with multiple format attempts
		var creationDatetime time.Time
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000000Z",
			"2006-01-02 15:04:05.000000",
		}

		parsed := false
		for _, format := range formats {
			if t, err := time.Parse(format, creationDatetimeStr); err == nil {
				creationDatetime = t
				parsed = true
				break
			}
		}

		if !parsed {
			// If parsing fails, use the original string
			item.Datetime = creationDatetimeStr
		} else {
			item.Datetime = creationDatetime.Format("2006-01-02 15:04:05")
		}

		sweatRate = append(sweatRate, item)
	}

	// Get sweat loss data
	sweatLossQuery := fmt.Sprintf(`
		SELECT creation_datetime, sweat_loss 
		FROM user_data 
		%s 
		ORDER BY creation_datetime DESC
	`, whereClause)

	sweatLossRows, err := s.db.Query(sweatLossQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sweat loss data: %v", err)
	}
	defer sweatLossRows.Close()

	var sweatLoss []models.HistoricalDataItem
	for sweatLossRows.Next() {
		var item models.HistoricalDataItem
		var creationDatetimeStr string
		err := sweatLossRows.Scan(&creationDatetimeStr, &item.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sweat loss data: %v", err)
		}

		// Parse the timestamp string with multiple format attempts
		var creationDatetime time.Time
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000000Z",
			"2006-01-02 15:04:05.000000",
		}

		parsed := false
		for _, format := range formats {
			if t, err := time.Parse(format, creationDatetimeStr); err == nil {
				creationDatetime = t
				parsed = true
				break
			}
		}

		if !parsed {
			// If parsing fails, use the original string
			item.Datetime = creationDatetimeStr
		} else {
			item.Datetime = creationDatetime.Format("2006-01-02 15:04:05")
		}

		sweatLoss = append(sweatLoss, item)
	}

	return &models.HistoricalDataResponse{
		SweatPosition: sweatPosition,
		SweatRate:     sweatRate,
		SweatLoss:     sweatLoss,
	}, nil
}
