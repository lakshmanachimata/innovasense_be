package services

import (
	"database/sql"
	"errors"
	"innovasense_be/config"
	"innovasense_be/models"
	"time"
)

type UserService struct {
	db *sql.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: config.GetDB(),
	}
}

// CheckUser validates user credentials
func (s *UserService) CheckUser(cnumber, userpin string) (*models.User, error) {
	query := `
		SELECT id, cnumber, userpin, username, gender, age, height, weight, 
		       role_id, ustatus, creation_datetime
		FROM users_master 
		WHERE cnumber = ? AND ustatus = 0
	`

	var user models.User
	var creationDatetimeStr string
	err := s.db.QueryRow(query, cnumber).Scan(
		&user.ID, &user.CNumber, &user.Userpin, &user.Username, &user.Gender,
		&user.Age, &user.Height, &user.Weight, &user.RoleID, &user.UStatus,
		&creationDatetimeStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Decrypt the stored userpin and compare with the input
	// encryptService := NewEncryptDecryptService()
	// fmt.Printf("Attempting to decrypt userpin: %s\n", user.Userpin)
	// decryptedUserpin, err := encryptService.GetDecryptData(user.Userpin)
	// if err != nil {
	// 	fmt.Printf("Decryption error: %v\n", err)
	// 	return nil, errors.New("invalid credentials")
	// }
	// fmt.Printf("Decrypted userpin: %s, Input userpin: %s\n", decryptedUserpin, userpin)

	// if userpin != decryptedUserpin {
	// 	return nil, errors.New("invalid credentials")
	// }

	// Parse the creation_datetime string to time.Time
	if creationDatetimeStr != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999", creationDatetimeStr)
		if err != nil {
			// Try alternative format without microseconds
			parsedTime, err := time.Parse("2006-01-02 15:04:05", creationDatetimeStr)
			if err != nil {
				// If parsing fails, set to zero time
				user.CreationDatetime = time.Time{}
			} else {
				user.CreationDatetime = parsedTime
			}
		} else {
			user.CreationDatetime = parsedTime
		}
	}

	return &user, nil
}

// CheckPIN validates user PIN
func (s *UserService) CheckPIN(id int, userpin string) (*models.User, error) {
	query := `
		SELECT id, cnumber, userpin, username, gender, age, height, weight, 
		       role_id, ustatus, creation_datetime
		FROM users_master 
		WHERE id = ? AND userpin = ?
	`

	var user models.User
	var creationDatetimeStr string
	err := s.db.QueryRow(query, id, userpin).Scan(
		&user.ID, &user.CNumber, &user.Userpin, &user.Username, &user.Gender,
		&user.Age, &user.Height, &user.Weight, &user.RoleID, &user.UStatus,
		&creationDatetimeStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid PIN")
		}
		return nil, err
	}

	// Parse the creation_datetime string to time.Time
	if creationDatetimeStr != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999", creationDatetimeStr)
		if err != nil {
			// Try alternative format without microseconds
			parsedTime, err = time.Parse("2006-01-02 15:04:05", creationDatetimeStr)
			if err != nil {
				// If parsing fails, set to zero time
				user.CreationDatetime = time.Time{}
			} else {
				user.CreationDatetime = parsedTime
			}
		} else {
			user.CreationDatetime = parsedTime
		}
	}

	return &user, nil
}

// ValidateUser checks if user exists and is active
func (s *UserService) ValidateUser(cnumber string) (*models.User, error) {
	query := `
		SELECT id, cnumber, userpin, username, gender, age, height, weight, 
		       role_id, ustatus, creation_datetime
		FROM users_master 
		WHERE cnumber = ? AND ustatus = 0
	`

	var user models.User
	var creationDatetimeStr string
	err := s.db.QueryRow(query, cnumber).Scan(
		&user.ID, &user.CNumber, &user.Userpin, &user.Username, &user.Gender,
		&user.Age, &user.Height, &user.Weight, &user.RoleID, &user.UStatus,
		&creationDatetimeStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Parse the creation_datetime string to time.Time
	if creationDatetimeStr != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05.999999", creationDatetimeStr)
		if err != nil {
			// Try alternative format without microseconds
			parsedTime, err = time.Parse("2006-01-02 15:04:05", creationDatetimeStr)
			if err != nil {
				// If parsing fails, set to zero time
				user.CreationDatetime = time.Time{}
			} else {
				user.CreationDatetime = parsedTime
			}
		} else {
			user.CreationDatetime = parsedTime
		}
	}

	return &user, nil
}

// AdminLogin validates admin credentials

// RegisterUser creates a new user
func (s *UserService) RegisterUser(req *models.RegisterRequest) (int, error) {
	query := `
		INSERT INTO users_master (cnumber, userpin, username, gender, age, 
		                         height, weight, creation_datetime, ustatus, role_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, 2)
	`

	now := time.Now()
	result, err := s.db.Exec(query, req.CNumber, req.Userpin, req.Username,
		req.Gender, req.Age, req.Height, req.Weight, now)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// ChangePassword updates user password
func (s *UserService) ChangePassword(userID int, oldPassword, newPassword string) error {
	// First verify the old password
	_, err := s.CheckPIN(userID, oldPassword)
	if err != nil {
		return errors.New("invalid old password")
	}

	// Update the password
	query := `UPDATE users_master SET userpin = ? WHERE id = ?`
	_, err = s.db.Exec(query, newPassword, userID)
	return err
}

// DeleteAccount marks user account as deleted
func (s *UserService) DeleteAccount(userID int) error {
	query := `UPDATE users_master SET ustatus = 5 WHERE id = ?`
	_, err := s.db.Exec(query, userID)
	return err
}

// GetData retrieves all data from a table (matches PHP getData method)
func (s *UserService) GetData(tableName string) ([]map[string]interface{}, error) {
	query := `SELECT * FROM ` + tableName

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		data = append(data, row)
	}

	return data, nil
}

// GetDataById retrieves specific fields from a table by ID (matches PHP getDataById)
func (s *UserService) GetDataById(tableName string, id int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, user_id, weight, height, sweat_position, time_taken, bmi, tbsa, 
		       sweat_rate, sweat_loss, creation_datetime, image_path
		FROM ` + tableName + ` 
		WHERE id = ?
	`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		data = append(data, row)
	}

	return data, nil
}

// GetNewDataById retrieves all data from a table by ID (matches PHP getNewDataById)
func (s *UserService) GetNewDataById(tableName string, id int) ([]map[string]interface{}, error) {
	query := `SELECT * FROM ` + tableName + ` WHERE id = ?`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		data = append(data, row)
	}

	return data, nil
}

// GlobalInsert inserts data into any table (matches PHP globalinsert)
func (s *UserService) GlobalInsert(tableName string, data map[string]interface{}) (int, error) {
	// Build dynamic INSERT query
	var columns []string
	var placeholders []string
	var values []interface{}

	for col, val := range data {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}

	query := `INSERT INTO ` + tableName + ` (` + columns[0]
	for i := 1; i < len(columns); i++ {
		query += `, ` + columns[i]
	}
	query += `) VALUES (` + placeholders[0]
	for i := 1; i < len(placeholders); i++ {
		query += `, ` + placeholders[i]
	}
	query += `)`

	result, err := s.db.Exec(query, values...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GlobalUpdate updates data in any table (matches PHP globalupdate)
func (s *UserService) GlobalUpdate(tableName string, data map[string]interface{}, id int) (int, error) {
	// Build dynamic UPDATE query
	var setClause []string
	var values []interface{}

	for col, val := range data {
		setClause = append(setClause, col+" = ?")
		values = append(values, val)
	}
	values = append(values, id)

	query := `UPDATE ` + tableName + ` SET ` + setClause[0]
	for i := 1; i < len(setClause); i++ {
		query += `, ` + setClause[i]
	}
	query += ` WHERE id = ?`

	_, err := s.db.Exec(query, values...)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUserIDByCNumber gets user ID by CNumber
func (s *UserService) GetUserIDByCNumber(cnumber string) (int, error) {
	query := `SELECT id FROM users_master WHERE cnumber = ? AND ustatus = 0`

	var userID int
	err := s.db.QueryRow(query, cnumber).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}
		return 0, err
	}

	return userID, nil
}
