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
	Type            string `json:"type" binding:"required"`
	Description     string `json:"description"`
	BackgroundImage string `json:"background_image"`
	LogoImage       string `json:"logo_image"`
	BodyTemplate    string `json:"body_template"`
	LayoutJSON      string `json:"layout_json"`
	PageSize        string `json:"page_size"`
	Orientation     string `json:"orientation"`
	FontFamily      string `json:"font_family"`
	FontSize        int    `json:"font_size"`
	FontColor       string `json:"font_color"`
	IsActive        *bool  `json:"is_active"`
}

type UpdateTemplateRequest struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	Description     string `json:"description"`
	BackgroundImage string `json:"background_image"`
	LogoImage       string `json:"logo_image"`
	BodyTemplate    string `json:"body_template"`
	LayoutJSON      string `json:"layout_json"`
	PageSize        string `json:"page_size"`
	Orientation     string `json:"orientation"`
	FontFamily      string `json:"font_family"`
	FontSize        int    `json:"font_size"`
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

	staffID, _ := c.Get("staffID")
	sid, _ := staffID.(uint)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	tmpl := models.CertificateTemplate{
		Name:            req.Name,
		Type:            req.Type,
		Description:     req.Description,
		BackgroundImage: req.BackgroundImage,
		LogoImage:       req.LogoImage,
		BodyTemplate:    req.BodyTemplate,
		LayoutJSON:      req.LayoutJSON,
		PageSize:        req.PageSize,
		Orientation:     req.Orientation,
		FontFamily:      req.FontFamily,
		FontSize:        req.FontSize,
		FontColor:       req.FontColor,
		IsActive:        isActive,
		CreatedByID:     &sid,
	}

	// Defaults
	if tmpl.FontFamily == "" {
		tmpl.FontFamily = "Arial"
	}
	if tmpl.FontSize == 0 {
		tmpl.FontSize = 16
	}
	if tmpl.FontColor == "" {
		tmpl.FontColor = "#000000"
	}
	if tmpl.PageSize == "" {
		tmpl.PageSize = "A4"
	}
	if tmpl.Orientation == "" {
		tmpl.Orientation = "landscape"
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

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.BackgroundImage != "" {
		updates["background_image"] = req.BackgroundImage
	}
	if req.LogoImage != "" {
		updates["logo_image"] = req.LogoImage
	}
	if req.BodyTemplate != "" {
		updates["body_template"] = req.BodyTemplate
	}
	if req.LayoutJSON != "" {
		updates["layout_json"] = req.LayoutJSON
	}
	if req.PageSize != "" {
		updates["page_size"] = req.PageSize
	}
	if req.Orientation != "" {
		updates["orientation"] = req.Orientation
	}
	if req.FontFamily != "" {
		updates["font_family"] = req.FontFamily
	}
	if req.FontSize > 0 {
		updates["font_size"] = req.FontSize
	}
	if req.FontColor != "" {
		updates["font_color"] = req.FontColor
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.db.Model(&tmpl).Updates(updates).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Shablonni yangilashda xatolik")
		return
	}

	h.db.First(&tmpl, id)
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
