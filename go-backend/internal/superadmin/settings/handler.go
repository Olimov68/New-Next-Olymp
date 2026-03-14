package superadminsettings

import (
	"net/http"

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

type UpdateGlobalSettingsRequest struct {
	PlatformName        *string `json:"platform_name"`
	DefaultLanguage     *string `json:"default_language"`
	SupportEmail        *string `json:"support_email"`
	MaintenanceMode     *bool   `json:"maintenance_mode"`
	RegistrationEnabled *bool   `json:"registration_enabled"`
}

// GetAll GET /api/v1/superadmin/settings
func (h *Handler) GetAll(c *gin.Context) {
	var setting models.GlobalSetting
	result := h.db.First(&setting)
	if result.Error != nil {
		// Default bilan yaratish
		setting = models.GlobalSetting{
			PlatformName:        "NextOlymp",
			DefaultLanguage:     "uz",
			SupportEmail:        "",
			MaintenanceMode:     false,
			RegistrationEnabled: true,
		}
		h.db.Create(&setting)
	}
	response.Success(c, http.StatusOK, "Global settings", setting)
}

// Update PUT /api/v1/superadmin/settings
func (h *Handler) Update(c *gin.Context) {
	var req UpdateGlobalSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Ensure a record exists
	var setting models.GlobalSetting
	if h.db.First(&setting).Error != nil {
		setting = models.GlobalSetting{
			PlatformName:        "NextOlymp",
			DefaultLanguage:     "uz",
			MaintenanceMode:     false,
			RegistrationEnabled: true,
		}
		h.db.Create(&setting)
	}

	fields := map[string]interface{}{}
	if req.PlatformName != nil {
		fields["platform_name"] = *req.PlatformName
	}
	if req.DefaultLanguage != nil {
		fields["default_language"] = *req.DefaultLanguage
	}
	if req.SupportEmail != nil {
		fields["support_email"] = *req.SupportEmail
	}
	if req.MaintenanceMode != nil {
		fields["maintenance_mode"] = *req.MaintenanceMode
	}
	if req.RegistrationEnabled != nil {
		fields["registration_enabled"] = *req.RegistrationEnabled
	}

	if len(fields) > 0 {
		h.db.Model(&setting).Updates(fields)
	}

	h.db.First(&setting, setting.ID)
	response.Success(c, http.StatusOK, "Settings updated", setting)
}
