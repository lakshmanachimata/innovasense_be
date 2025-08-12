package controllers

import (
	"innovasense_be/models"
	"innovasense_be/services"
	"net/http"

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
// @Tags Data Retrieval
// @Accept json
// @Produce json
// @Param request body object{} true "Empty request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/getSweatImages [post]
func (c *CommonController) GetSweatImages(ctx *gin.Context) {
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
// @Description Upload an image and save the path
// @Tags File Management
// @Accept json
// @Produce json
// @Param request body models.ImageUploadRequest true "Image upload request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/uploadInnovoImage [post]
func (c *CommonController) UploadInnovoImage(ctx *gin.Context) {
	var req models.ImageUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
		})
		return
	}

	id, err := c.commonService.SaveImagePath(req.UserID, req.ImagePath)
	if err != nil {
		ctx.JSON(http.StatusOK, models.APIResponse{
			Code:     1,
			Message:  "Failed to upload image",
			Response: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Code:     0,
		Message:  "Image uploaded successfully",
		Response: id,
	})
}

// UpdateInnovoImagePath handles image path update
// @Summary Update image path
// @Description Update the path of an existing image
// @Tags File Management
// @Accept json
// @Produce json
// @Param request body models.UpdateImagePathRequest true "Image path update request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /Services/protected/updateInnovoImagePath [post]
func (c *CommonController) UpdateInnovoImagePath(ctx *gin.Context) {
	var req models.UpdateImagePathRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    1,
			Message: "Invalid request data",
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
