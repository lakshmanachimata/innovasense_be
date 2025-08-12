package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"innovasense_be/models"
	"innovasense_be/services"

	"innovasense_be/middleware"

	"github.com/gin-gonic/gin"
)

type CommonController struct {
	commonService *services.CommonService
	userService   *services.UserService
}

func NewCommonController() *CommonController {
	return &CommonController{
		commonService: services.NewCommonService(),
		userService:   services.NewUserService(),
	}
}

// GetBannerImages handles banner images retrieval
// @Summary Get banner images
// @Description Retrieve banner images
// @Tags Data Retrieval
// @Accept json
// @Produce json
// @Param request body object{} true "Empty request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/getBannerImages [post]
func (c *CommonController) GetBannerImages(ctx *gin.Context) {
	images, err := c.commonService.GetBannerImages()
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get banner images",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: images,
	})
}

// GetHomeImages handles home images retrieval
// @Summary Get home images
// @Description Retrieve home page images
// @Tags Data Retrieval
// @Accept json
// @Produce json
// @Param request body object{} true "Empty request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/getHomeImages [post]
func (c *CommonController) GetHomeImages(ctx *gin.Context) {
	images, err := c.commonService.GetHomeImages()
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get home images",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: images,
	})
}

// GetSweatImages handles sweat images retrieval
// @Summary Get sweat images
// @Description Retrieve sweat analysis images with metadata
// @Tags File Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.GetSweatImagesRequest true "Sweat images request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getSweatImages [post]
func (c *CommonController) GetSweatImages(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.GetSweatImagesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Validate cnumber and username from request body against JWT claims
	if req.CNumber != claims.CNumber {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "cnumber in request body does not match authenticated user",
		})
		return
	}

	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	images, err := c.commonService.GetSweatImages()
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get sweat images",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: images,
	})
}

// GetDevices handles devices retrieval
// @Summary Get devices
// @Description Retrieve available device types
// @Tags Data Retrieval
// @Accept json
// @Produce json
// @Param request body object{} true "Empty request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/getDevices [post]
func (c *CommonController) GetDevices(ctx *gin.Context) {
	devices, err := c.commonService.GetDevices()
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to get devices",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "OK",
		Response: devices,
	})
}

// UploadInnovoImage handles image upload
// @Summary Upload image
// @Description Upload an image and save it to server with username_timestamp.jpg format
// @Tags File Management
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.ImageUploadRequest true "User identity information"
// @Param image formData file true "Image file to upload"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/uploadInnovoImage [post]
func (c *CommonController) UploadInnovoImage(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	// Parse multipart form data
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Failed to parse form data",
		})
		return
	}

	// Get JSON request body from form field
	requestBody := ctx.PostForm("request")
	if requestBody == "" {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Request body is required",
		})
		return
	}

	// Parse the JSON request body
	var req models.ImageUploadRequest
	if err := json.Unmarshal([]byte(requestBody), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request body format",
		})
		return
	}

	// Validate cnumber and username against JWT claims
	if req.CNumber != claims.CNumber {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "cnumber in request does not match authenticated user",
		})
		return
	}

	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request does not match authenticated user",
		})
		return
	}

	// Get user ID from CNumber using user service
	userID, err := c.userService.GetUserIDByCNumber(req.CNumber)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "User not found",
			Response: 0,
		})
		return
	}

	// Get uploaded file
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "No image file uploaded",
		})
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidImageType(header.Filename) {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid file type. Only JPG, JPEG, and PNG are allowed",
		})
		return
	}

	// Create assets/innovo directory if it doesn't exist
	uploadDir := "assets/innovo"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    1,
			Message: "Failed to create upload directory",
		})
		return
	}

	// Generate filename: username_timestamp.jpg
	timestamp := time.Now().Format("20060102_150405")
	fileExt := getFileExtension(header.Filename)
	filename := fmt.Sprintf("%s_%s%s", req.Username, timestamp, fileExt)
	filepath := filepath.Join(uploadDir, filename)

	// Create the file on disk
	dst, err := os.Create(filepath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    1,
			Message: "Failed to create file on server",
		})
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    1,
			Message: "Failed to save uploaded file",
		})
		return
	}

	// Save image path to database
	id, err := c.commonService.SaveImagePath(userID, filepath)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to save image path to database",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "Image uploaded successfully",
		Response: map[string]interface{}{
			"id":          id,
			"filename":    filename,
			"filepath":    filepath,
			"size":        header.Size,
			"uploaded_at": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// UpdateInnovoImagePath handles image path update
// @Summary Update image path
// @Description Update the path of an existing image
// @Tags File Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param request body models.UpdateImagePathRequest true "Image path update request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/updateInnovoImagePath [post]
func (c *CommonController) UpdateInnovoImagePath(ctx *gin.Context) {
	// Get user information from JWT claims
	claims, exists := middleware.GetJWTClaimsFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    1,
			Message: "User not authenticated",
		})
		return
	}

	var req models.UpdateImagePathRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	// Validate cnumber and username from request body against JWT claims
	if req.CNumber != claims.CNumber {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "cnumber in request body does not match authenticated user",
		})
		return
	}

	if req.Username != claims.UserName {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Code:    1,
			Message: "username in request body does not match authenticated user",
		})
		return
	}

	err := c.commonService.UpdateImagePath(req.UserID, req.ImageID, req.ImagePath)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to update image path",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "Image path updated successfully",
	})
}

// Helper function to validate image file types
func isValidImageType(filename string) bool {
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}

// Helper function to get file extension
func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".jpg" // Default to .jpg if no extension
	}
	return ext
}
