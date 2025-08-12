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
