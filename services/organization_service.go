package services

import (
	"database/sql"
	"errors"
	"fmt"
	"innovasense_be/models"
)

type OrganizationService struct {
	db *sql.DB
}

func NewOrganizationService(db *sql.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

// ValidateOrgCredentials validates API key and secret key to get organization ID
func (s *OrganizationService) ValidateOrgCredentials(apiKey, secretKey string) (*models.Organization, error) {
	query := `
		SELECT id, org_name, org_desc, salt_key, api_key 
		FROM organizations 
		WHERE api_key = ? AND salt_key = ?
	`

	fmt.Printf("DEBUG: Executing query: %s\n", query)
	fmt.Printf("DEBUG: With params: apiKey=%s, secretKey=%s\n", apiKey, secretKey)

	var org models.Organization
	err := s.db.QueryRow(query, apiKey, secretKey).Scan(
		&org.ID,
		&org.OrgName,
		&org.OrgDesc,
		&org.SaltKey,
		&org.APIKey,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid API key or secret key")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &org, nil
}

// CheckUserExists checks if a user exists by contact number or email for a specific organization
func (s *OrganizationService) CheckUserExists(contact, email string, orgID int) (*models.OrgUser, error) {
	query := `
		SELECT id, email_id, user_pwd, user_name, org_id 
		FROM org_users 
		WHERE (email_id = ? OR email_id = ?) AND org_id = ?
	`

	var user models.OrgUser
	err := s.db.QueryRow(query, contact, email, orgID).Scan(
		&user.ID,
		&user.EmailID,
		&user.UserPwd,
		&user.UserName,
		&user.OrgID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User doesn't exist
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &user, nil
}

// RegisterOrgUser registers a new user in the organization
func (s *OrganizationService) RegisterOrgUser(name, contact string, orgID int) (*models.OrgUser, error) {
	// Check if contact is email or phone number
	// For simplicity, we'll treat it as email if it contains @, otherwise as phone
	emailID := contact

	query := `
		INSERT INTO org_users (email_id, user_pwd, user_name, org_id) 
		VALUES (?, ?, ?, ?)
	`

	result, err := s.db.Exec(query, emailID, "default123", name, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	// For MySQL, we can use LastInsertId
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return &models.OrgUser{
		ID:       int(userID),
		EmailID:  emailID,
		UserPwd:  "default123",
		UserName: name,
		OrgID:    orgID,
	}, nil
}

// GetUserIDByContact gets user ID from users_master table by contact number
func (s *OrganizationService) GetUserIDByContact(contact string) (int, error) {
	query := `SELECT id FROM users_master WHERE cnumber = ?`

	var userID int
	err := s.db.QueryRow(query, contact).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found in users_master")
		}
		return 0, fmt.Errorf("database error: %v", err)
	}

	return userID, nil
}

// CreateUserInMaster creates a new user in users_master table
func (s *OrganizationService) CreateUserInMaster(name, contact, gender string, age int, height, weight float64, orgID int) (int, error) {
	query := `
		INSERT INTO users_master (cnumber, username, gender, age, height, weight, ustatus, role_id) 
		VALUES (?, ?, ?, ?, ?, ?, 0, 2)
	`

	result, err := s.db.Exec(query, contact, name, gender, age, height, weight)
	if err != nil {
		return 0, fmt.Errorf("failed to create user in master: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID: %v", err)
	}

	return int(userID), nil
}

// CheckUserExistsByContact checks if a user exists by contact for a specific organization
func (s *OrganizationService) CheckUserExistsByContact(contact string, orgID int) (*models.OrgUser, error) {
	query := `
		SELECT id, email_id, user_pwd, user_name, org_id 
		FROM org_users 
		WHERE email_id = ? AND org_id = ?
	`

	var user models.OrgUser
	err := s.db.QueryRow(query, contact, orgID).Scan(
		&user.ID,
		&user.EmailID,
		&user.UserPwd,
		&user.UserName,
		&user.OrgID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User doesn't exist
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &user, nil
}
