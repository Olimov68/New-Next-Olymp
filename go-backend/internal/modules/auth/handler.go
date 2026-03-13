package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/utils"
	"github.com/nextolympservice/go-backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body RegisterRequest true "Register request"
// @Success      201 {object} response.Response{data=RegisterResponse}
// @Failure      422 {object} response.Response
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	result, err := h.service.Register(&req)
	if err != nil {
		// Biznes logika xatoliklari 400 qaytaradi
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Registration successful", result)
}

// Login godoc
// @Summary      Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body LoginRequest true "Login request"
// @Success      200 {object} response.Response{data=LoginResponse}
// @Failure      401 {object} response.Response
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	result, err := h.service.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", result)
}

// Logout godoc
// @Summary      Logout user (client-side token discard)
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	// Sodda variant: client tokenni o'chiradi
	// Kelajakda token blacklist qo'shish mumkin
	response.Success(c, http.StatusOK, "Logged out successfully", nil)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body RefreshRequest true "Refresh token"
// @Success      200 {object} response.Response{data=TokenPair}
// @Failure      401 {object} response.Response
// @Router       /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	tokens, err := h.service.RefreshTokens(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed", tokens)
}

// Me godoc
// @Summary      Get current user info
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=MeResponse}
// @Failure      401 {object} response.Response
// @Router       /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	result, err := h.service.GetMe(userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User info", result)
}
