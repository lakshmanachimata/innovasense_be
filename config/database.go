package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	// Check database type from environment
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "mysql" // Default to MySQL
	}

	var err error

	if dbType == "postgres" {
		// Use PostgreSQL database
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}

		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}

		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "lakshmana"
		}

		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			dbPassword = ""
		}

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "innosense"
		}

		// Create PostgreSQL connection string
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbPort, dbName)

		// Open database connection
		DB, err = sql.Open("postgres", dsn)
		if err != nil {
			return fmt.Errorf("failed to open PostgreSQL database: %v", err)
		}

		// Test the connection
		if err := DB.Ping(); err != nil {
			return fmt.Errorf("failed to ping PostgreSQL database: %v", err)
		}

		// Set connection pool settings
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)

		log.Println("PostgreSQL database connected successfully")
	} else if dbType == "sqlite" {
		// Use in-memory SQLite database for testing
		DB, err = sql.Open("sqlite3", ":memory:")
		if err != nil {
			return fmt.Errorf("failed to open SQLite database: %v", err)
		}

		// Enable foreign keys
		_, err = DB.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			return fmt.Errorf("failed to enable foreign keys: %v", err)
		}

		// Create tables
		if err := createTables(); err != nil {
			return fmt.Errorf("failed to create tables: %v", err)
		}

		log.Println("SQLite in-memory database connected successfully")
		log.Println("Database tables created and sample data inserted")
	} else {
		// Use MySQL database
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}

		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "3306"
		}

		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "root"
		}

		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			dbPassword = ""
		}

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "innosense"
		}

		// Create database connection string
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName)

		// Open database connection
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			return fmt.Errorf("failed to open database: %v", err)
		}

		log.Printf("DEBUG: Connected to MySQL database with DSN: %s", dsn)

		// Test the connection
		if err := DB.Ping(); err != nil {
			return fmt.Errorf("failed to ping database: %v", err)
		}

		// Set connection pool settings
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)

		log.Println("MySQL database connected successfully")
	}

	return nil
}

func GetDB() *sql.DB {
	return DB
}

// createTables creates all necessary tables for the application
func createTables() error {
	tables := []string{
		// Users master table
		`CREATE TABLE IF NOT EXISTS users_master (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			cnumber TEXT,
			userpin TEXT NOT NULL,
			username TEXT NOT NULL,
			gender TEXT,
			age INTEGER,
			height REAL,
			weight REAL,
			role_id INTEGER DEFAULT 2,
			ustatus INTEGER DEFAULT 0,
			creation_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// User data table for hydration records
		`CREATE TABLE IF NOT EXISTS user_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			weight REAL,
			height REAL,
			sweat_position REAL,
			time_taken REAL,
			bmi REAL,
			tbsa REAL,
			image_path TEXT,
			sweat_rate REAL,
			sweat_loss REAL,
			device_type INTEGER,
			image_id INTEGER,
			creation_datetime DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users_master(id)
		)`,

		// Sweat data table
		`CREATE TABLE IF NOT EXISTS sweat_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			image_id INTEGER NOT NULL,
			sweat_rate REAL,
			sweat_loss REAL,
			creation_datetime DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users_master(id)
		)`,

		// Sweat summary table
		`CREATE TABLE IF NOT EXISTS sweat_summary (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			low_limit REAL,
			high_limit REAL,
			hyd_status TEXT,
			comments TEXT,
			recomm TEXT,
			color TEXT
		)`,

		// Sweat rate summary table
		`CREATE TABLE IF NOT EXISTS sweat_rate_summary (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			low_limit REAL,
			high_limit REAL,
			hyd_status TEXT,
			comments TEXT,
			recomm TEXT,
			color TEXT
		)`,

		// Sweat images table
		`CREATE TABLE IF NOT EXISTS sweat_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			image_path TEXT NOT NULL,
			sweat_range TEXT,
			implications TEXT,
			recomm TEXT,
			strategy TEXT,
			result TEXT,
			colorcode TEXT
		)`,

		// Banner images table
		`CREATE TABLE IF NOT EXISTS banner_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			image_path TEXT NOT NULL
		)`,

		// Home images table
		`CREATE TABLE IF NOT EXISTS home_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			image_path TEXT NOT NULL
		)`,

		// Device master table
		`CREATE TABLE IF NOT EXISTS device_master (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_name TEXT NOT NULL,
			device_text TEXT
		)`,

		// Organization table
		`CREATE TABLE IF NOT EXISTS organizations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			org_name TEXT NOT NULL,
			org_desc TEXT,
			salt_key TEXT,
			api_key TEXT UNIQUE
		)`,

		// Organization users table
		`CREATE TABLE IF NOT EXISTS org_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email_id TEXT NOT NULL,
			user_pwd TEXT NOT NULL,
			user_name TEXT NOT NULL,
			org_id INTEGER NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id)
		)`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	// Insert sample data
	if err := insertSampleData(); err != nil {
		return fmt.Errorf("failed to insert sample data: %v", err)
	}

	log.Println("All tables created successfully")
	return nil
}

// insertSampleData inserts sample data for testing
func insertSampleData() error {
	// Insert banner images
	bannerImages := []string{
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/1.png')",
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/2.png')",
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/3.jpg')",
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/4.jpg')",
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/5.jpg')",
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/banner1.png')",
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/banner2.png')",
		"INSERT INTO banner_images (image_path) VALUES ('/assets/banners/banner3.png')",
	}

	for _, query := range bannerImages {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Warning: Failed to insert banner image: %v", err)
		}
	}

	// Insert home images (same as banner images for now)
	homeImages := []string{
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/1.png')",
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/2.png')",
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/3.jpg')",
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/4.jpg')",
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/5.jpg')",
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/banner1.png')",
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/banner2.png')",
		"INSERT INTO home_images (image_path) VALUES ('/assets/banners/banner3.png')",
	}

	for _, query := range homeImages {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Warning: Failed to insert home image: %v", err)
		}
	}

	// Insert devices
	devices := []string{
		"INSERT INTO device_master (device_name, device_text) VALUES ('Hydrosense (Classic)', 'Standard sweat hydration tracking for average exercise and sweat loss.')",
		"INSERT INTO device_master (device_name, device_text) VALUES ('Hydrosense Plus (+)', 'High-volume hydration tracking for longer or more intense workouts.')",
		"INSERT INTO device_master (device_name, device_text) VALUES ('Hydrosense Pro', 'Sweat hydration + electrolyte monitoring for regular training sessions.')",
		"INSERT INTO device_master (device_name, device_text) VALUES ('Hydrosense Pro Plus (+)', 'Full hydration and electrolyte analysis for high-intensity or extended activity.')",
	}

	for _, query := range devices {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Warning: Failed to insert device: %v", err)
		}
	}

	// Insert sample organization
	orgQuery := `INSERT INTO organizations (org_name, org_desc, salt_key, api_key) 
		VALUES ('Test Organization', 'Test organization for API testing', 'test-salt-key', 'test-api-key')`
	if _, err := DB.Exec(orgQuery); err != nil {
		log.Printf("Warning: Failed to insert organization: %v", err)
	}

	// Insert sample organization user
	orgUserQuery := `INSERT INTO org_users (email_id, user_pwd, user_name, org_id) 
		VALUES ('test@example.com', 'test123', 'Test User', 1)`
	if _, err := DB.Exec(orgUserQuery); err != nil {
		log.Printf("Warning: Failed to insert organization user: %v", err)
	}

	// Insert sample sweat summary data
	sweatSummaryQueries := []string{
		"INSERT INTO sweat_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES (0.0, 0.5, 'Low', 'Low hydration level', 'Increase fluid intake', 'red')",
		"INSERT INTO sweat_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES (0.5, 1.0, 'Normal', 'Normal hydration level', 'Maintain current fluid intake', 'green')",
		"INSERT INTO sweat_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES (1.0, 2.0, 'High', 'High hydration level', 'Consider reducing fluid intake', 'yellow')",
	}

	for _, query := range sweatSummaryQueries {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Warning: Failed to insert sweat summary: %v", err)
		}
	}

	// Insert sample sweat rate summary data
	sweatRateSummaryQueries := []string{
		"INSERT INTO sweat_rate_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES (0.0, 0.5, 'Low Rate', 'Low sweat rate', 'Increase activity intensity', 'blue')",
		"INSERT INTO sweat_rate_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES (0.5, 1.5, 'Normal Rate', 'Normal sweat rate', 'Maintain current activity level', 'green')",
		"INSERT INTO sweat_rate_summary (low_limit, high_limit, hyd_status, comments, recomm, color) VALUES (1.5, 3.0, 'High Rate', 'High sweat rate', 'Consider reducing intensity', 'orange')",
	}

	for _, query := range sweatRateSummaryQueries {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Warning: Failed to insert sweat rate summary: %v", err)
		}
	}

	log.Println("Sample data inserted successfully")
	return nil
}
