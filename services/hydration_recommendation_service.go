package services

import (
	"fmt"
	"innovasense_be/models"
)

type HydrationRecommendationService struct {
	hydrationService *HydrationService
	orgService       *OrganizationService
	userService      *UserService
}

func NewHydrationRecommendationService(
	hydrationService *HydrationService,
	orgService *OrganizationService,
	userService *UserService,
) *HydrationRecommendationService {
	return &HydrationRecommendationService{
		hydrationService: hydrationService,
		orgService:       orgService,
		userService:      userService,
	}
}

// GetHydrationRecommendation processes the hydration recommendation request
func (s *HydrationRecommendationService) GetHydrationRecommendation(
	req *models.HydrationRecommendationRequest,
	apiKey, secretKey string,
) (*models.EnhancedHydrationResponse, error) {

	// Step 1: Validate organization credentials
	org, err := s.orgService.ValidateOrgCredentials(apiKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("organization validation failed: %v", err)
	}

	// Step 2: Check if user exists in org_users table
	orgUser, err := s.orgService.CheckUserExists(req.Contact, req.Contact, org.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %v", err)
	}

	var userID int

	if orgUser != nil {
		// User exists in org_users, get user ID from users_master
		userID, err = s.orgService.GetUserIDByContact(req.Contact)
		if err != nil {
			// User exists in org_users but not in users_master, create in users_master
			userID, err = s.orgService.CreateUserInMaster(
				req.Name, req.Contact, req.Gender, req.Age, req.Height, req.Weight, org.ID,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create user in master: %v", err)
			}
		}
	} else {
		// User doesn't exist, register in both tables
		orgUser, err = s.orgService.RegisterOrgUser(req.Name, req.Contact, org.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to register user in organization: %v", err)
		}

		userID, err = s.orgService.CreateUserInMaster(
			req.Name, req.Contact, req.Gender, req.Age, req.Height, req.Weight, org.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create user in master: %v", err)
		}
	}

	// Step 3: Create hydration request and insert into user_data table
	hydrationReq := &models.HydrationRequest{
		CNumber:       req.Contact,
		Username:      req.Name,
		UserID:        userID,
		Weight:        req.Weight,
		Height:        req.Height,
		SweatPosition: req.SweatPosition,
		TimeTaken:     req.WorkoutTime,
		DeviceType:    1, // Default device type
	}

	// Insert hydration data and get the same response as existing insert service
	response, err := s.hydrationService.SaveEnhancedHydrationData(hydrationReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create hydration record: %v", err)
	}

	return response, nil
}
