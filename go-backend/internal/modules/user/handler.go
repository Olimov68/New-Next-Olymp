package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/modules/auth"
	"github.com/nextolympservice/go-backend/internal/utils"
	"github.com/nextolympservice/go-backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CompleteProfile godoc
// @Summary      Complete user profile (1st time)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CompleteProfileRequest true "Profile data"
// @Success      200 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /profile/complete [post]
func (h *Handler) CompleteProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	var req CompleteProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	// Photo file (optional in FormData)
	photoFile, _ := c.FormFile("photo")

	profile, err := h.service.CompleteProfile(userID.(uint), &req, photoFile)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Profile completed successfully", auth.ToProfileResponse(profile))
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body UpdateProfileRequest true "Profile data"
// @Success      200 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /profile/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	profile, err := h.service.UpdateProfile(userID.(uint), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", auth.ToProfileResponse(profile))
}

// UploadPhoto godoc
// @Summary      Upload profile photo
// @Tags         profile
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        photo formData file true "Photo file (jpg, png, webp)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /profile/photo [post]
func (h *Handler) UploadPhoto(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Photo file is required", nil)
		return
	}

	photoURL, err := h.service.UploadPhoto(userID.(uint), file)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Photo uploaded successfully", gin.H{
		"photo_url": photoURL,
	})
}

// GetProfile godoc
// @Summary      Get current user profile
// @Tags         profile
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Router       /profile/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	profile, err := h.service.GetProfile(userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile data", auth.ToProfileResponse(profile))
}
