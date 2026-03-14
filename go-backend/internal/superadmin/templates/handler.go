package templates

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type CreateTemplateRequest struct {
	Name            string `json:"name" binding:"required,min=2,max=200"`
	Type            string `json:"type" binding:"required,oneof=olympiad_winner olympiad_participant mock_test"`
	BackgroundImage string `json:"background_image"`
	LogoImage       string `json:"logo_image"`
	BodyTemplate    string `json:"body_template" binding:"required,min=10"`
	FontFamily      string `json:"font_family"`
	FontSize        int    `json:"font_size" binding:"omitempty,min=8,max=72"`
	FontColor       string `json:"font_color"`
}

type UpdateTemplateRequest struct {
	Name            string `json:"name" binding:"omitempty,min=2,max=200"`
	Type            string `json:"type" binding:"omitempty,oneof=olympiad_winner olympiad_participant mock_test"`
	BackgroundImage string `json:"background_image"`
	LogoImage       string `json:"logo_image"`
	BodyTemplate    string `json:"body_template" binding:"omitempty,min=10"`
	FontFamily      string `json:"font_family"`
	FontSize        int    `json:"font_size" binding:"omitempty,min=8,max=72"`
	FontColor       string `json:"font_color"`
	IsActive        *bool  `json:"is_active"`
}

// List — shablonlar ro'yxati
func (h *Handler) List(c *gin.Context) {
	var templates []models.CertificateTemplate
	query := h.db.Model(&models.CertificateTemplate{})

	tType := c.Query("type")
	if tType != "" {
		query = query.Where("type = ?", tType)
	}

	query.Order("created_at DESC").Find(&templates)
	response.Success(c, http.StatusOK, "Shablonlar", templates)
}

// GetByID
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	var tmpl models.CertificateTemplate
	if err := h.db.First(&tmpl, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Shablon topilmadi")
		return
	}

	response.Success(c, http.StatusOK, "Shablon", tmpl)
}

// Create
func (h *Handler) Create(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	staffID := c.GetUint("staff_id")
	tmpl := models.CertificateTemplate{
		Name:            req.Name,
		Type:            req.Type,
		BackgroundImage: req.BackgroundImage,
		LogoImage:       req.LogoImage,
		BodyTemplate:    req.BodyTemplate,
		FontFamily:      req.FontFamily,
		FontSize:        req.FontSize,
		FontColor:       req.FontColor,
		IsActive:        true,
		CreatedByID:     &staffID,
	}

	if tmpl.FontFamily == "" {
		tmpl.FontFamily = "Arial"
	}
	if tmpl.FontSize == 0 {
		tmpl.FontSize = 16
	}
	if tmpl.FontColor == "" {
		tmpl.FontColor = "#000000"
	}

	if err := h.db.Create(&tmpl).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Shablon yaratishda xatolik")
		return
	}

	response.Success(c, http.StatusCreated, "Shablon yaratildi", tmpl)
}

// Update
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	var tmpl models.CertificateTemplate
	if err := h.db.First(&tmpl, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Shablon topilmadi")
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if req.Name != "" {
		tmpl.Name = req.Name
	}
	if req.Type != "" {
		tmpl.Type = req.Type
	}
	if req.BackgroundImage != "" {
		tmpl.BackgroundImage = req.BackgroundImage
	}
	if req.LogoImage != "" {
		tmpl.LogoImage = req.LogoImage
	}
	if req.BodyTemplate != "" {
		tmpl.BodyTemplate = req.BodyTemplate
	}
	if req.FontFamily != "" {
		tmpl.FontFamily = req.FontFamily
	}
	if req.FontSize > 0 {
		tmpl.FontSize = req.FontSize
	}
	if req.FontColor != "" {
		tmpl.FontColor = req.FontColor
	}
	if req.IsActive != nil {
		tmpl.IsActive = *req.IsActive
	}

	if err := h.db.Save(&tmpl).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Shablonni yangilashda xatolik")
		return
	}

	response.Success(c, http.StatusOK, "Shablon yangilandi", tmpl)
}

// Delete
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	if err := h.db.Delete(&models.CertificateTemplate{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Shablonni o'chirishda xatolik")
		return
	}

	response.Success(c, http.StatusOK, "Shablon o'chirildi", nil)
}
