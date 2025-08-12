package services

import (
	"database/sql"
	"innovasense_be/config"
	"innovasense_be/models"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CommonService struct {
	db *sql.DB
}

func NewCommonService() *CommonService {
	return &CommonService{
		db: config.GetDB(),
	}
}

// GetBannerImages retrieves banner images from local file system
func (s *CommonService) GetBannerImages() ([]models.BannerImage, error) {
	bannersDir := "assets/banners"

	// Check if directory exists
	if _, err := os.Stat(bannersDir); os.IsNotExist(err) {
		return nil, err
	}

	// Read directory contents
	files, err := os.ReadDir(bannersDir)
	if err != nil {
		return nil, err
	}

	var images []models.BannerImage
	id := 1

	// Filter for image files and create banner image objects
	for _, file := range files {
		if !file.IsDir() {
			// Check if file is an image by extension
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" {
				// Create full accessible URL
				imageURL := "/assets/banners/" + file.Name()

				// Create banner image object
				image := models.BannerImage{
					ID:        id,
					ImagePath: imageURL,
				}

				images = append(images, image)
				id++
			}
		}
	}

	// Sort by filename for consistent ordering
	sort.Slice(images, func(i, j int) bool {
		return filepath.Base(images[i].ImagePath) < filepath.Base(images[j].ImagePath)
	})

	return images, nil
}

// GetHomeImages retrieves home images from local file system
func (s *CommonService) GetHomeImages() ([]models.HomeImage, error) {
	// Use banners directory for home images as well, or create a separate home directory
	homeDir := "assets/banners"

	// Check if directory exists
	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		return nil, err
	}

	// Read directory contents
	files, err := os.ReadDir(homeDir)
	if err != nil {
		return nil, err
	}

	var images []models.HomeImage
	id := 1

	// Filter for image files and create home image objects
	for _, file := range files {
		if !file.IsDir() {
			// Check if file is an image by extension
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" {
				// Create full accessible URL
				imageURL := "/assets/banners/" + file.Name()

				// Create home image object
				image := models.HomeImage{
					ID:        id,
					ImagePath: imageURL,
				}

				images = append(images, image)
				id++
			}
		}
	}

	// Sort by filename for consistent ordering
	sort.Slice(images, func(i, j int) bool {
		return filepath.Base(images[i].ImagePath) < filepath.Base(images[j].ImagePath)
	})

	return images, nil
}

// GetSweatImages retrieves sweat analysis images (matches PHP logic with id > 0 filter)
func (s *CommonService) GetSweatImages() ([]models.SweatImage, error) {
	query := `
		SELECT id, image_path, sweat_range, implications, recomm, strategy, result, colorcode
		FROM sweat_images 
		WHERE id > 0
		ORDER BY id
	`

	rows, err := s.db.Query(query)
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

// GetDevices retrieves device master data
func (s *CommonService) GetDevices() ([]models.Device, error) {
	query := `SELECT id, device_name, device_text FROM device_master ORDER BY id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(&device.ID, &device.DeviceName, &device.DeviceText)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// SaveImagePath saves image path to database
func (s *CommonService) SaveImagePath(userID int, imagePath string) (int, error) {
	query := `INSERT INTO user_images (user_id, image_path) VALUES (?, ?)`

	result, err := s.db.Exec(query, userID, imagePath)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// UpdateImagePath updates image path in database
func (s *CommonService) UpdateImagePath(userID, imageID int, imagePath string) error {
	query := `UPDATE user_images SET image_path = ? WHERE id = ? AND user_id = ?`

	_, err := s.db.Exec(query, imagePath, imageID, userID)
	return err
}

// GetClientHistory retrieves client history (matches PHP getClientHistory)
func (s *CommonService) GetClientHistory(userID int) ([]models.HydrationData, error) {
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
		err := rows.Scan(&item.ID, &item.UserID, &item.Weight, &item.Height, &item.SweatPosition,
			&item.TimeTaken, &item.BMI, &item.TBSA, &item.ImagePath, &item.SweatRate,
			&item.SweatLoss, &item.DeviceType, &item.ImageID, &item.CreationDatetime)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}
