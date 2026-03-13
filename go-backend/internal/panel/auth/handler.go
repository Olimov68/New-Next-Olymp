package panelauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/middleware"
	"github.com/nextolympservice/go-backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login — admin yoki superadmin login
// POST /api/v1/panel/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	res, err := h.service.Login(&req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}

// RefreshToken — token yangilash
// POST /api/v1/panel/auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	tokens, err := h.service.RefreshTokens(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed", tokens)
}

// Me — joriy panel foydalanuvchi ma'lumotlari
// GET /api/v1/panel/auth/me
func (h *Handler) Me(c *gin.Context) {
	staffID, exists := c.Get("staffID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	res, err := h.service.GetMe(staffID.(uint))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Staff info", res)
}

// Permissions — joriy foydalanuvchi ruxsatlari
// GET /api/v1/panel/auth/permissions
func (h *Handler) Permissions(c *gin.Context) {
	staffID, exists := c.Get("staffID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	role, _ := c.Get("staffRole")
	roleStr, _ := role.(string)

	permissions := middleware.GetStaffPermissionCodes(h.service.repo.GetDB(), staffID.(uint), roleStr)

	response.Success(c, http.StatusOK, "Permissions", gin.H{
		"permissions": permissions,
	})
}

// Logout — panel logout (frontend tokenni o'chiradi)
// POST /api/v1/panel/auth/logout
func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, http.StatusOK, "Logged out", nil)
}
