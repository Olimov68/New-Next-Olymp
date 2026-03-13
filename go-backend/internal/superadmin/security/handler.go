package sasecurity

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

type UpdateSecurityRequest struct {
	FullscreenRequired    *bool    `json:"fullscreen_required"`
	TabSwitchLimit        *int     `json:"tab_switch_limit"`
	BlurLimit             *int     `json:"blur_limit"`
	OfflineLimit          *int     `json:"offline_limit"`
	HeartbeatIntervalSecs *int     `json:"heartbeat_interval_seconds"`
	OneDevicePerSession   *bool    `json:"one_device_per_session"`
	ReconnectPolicy       *string  `json:"reconnect_policy"`
	RiskThreshold         *float64 `json:"risk_threshold"`
}

// GetSettings GET /api/v1/superadmin/security/settings
func (h *Handler) GetSettings(c *gin.Context) {
	var setting models.SecuritySetting
	result := h.db.First(&setting)
	if result.Error != nil {
		// Agar hali mavjud bo'lmasa, default bilan yaratamiz
		setting = models.SecuritySetting{
			FullscreenRequired:    true,
			TabSwitchLimit:        3,
			BlurLimit:             3,
			OfflineLimit:          2,
			HeartbeatIntervalSecs: 30,
			OneDevicePerSession:   true,
			ReconnectPolicy:       "allow_once",
			RiskThreshold:         0.7,
		}
		h.db.Create(&setting)
	}
	response.Success(c, http.StatusOK, "Security settings", setting)
}

// UpdateSettings PUT /api/v1/superadmin/security/settings
func (h *Handler) UpdateSettings(c *gin.Context) {
	var req UpdateSecurityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Ensure a record exists
	var setting models.SecuritySetting
	if h.db.First(&setting).Error != nil {
		setting = models.SecuritySetting{
			FullscreenRequired:    true,
			TabSwitchLimit:        3,
			BlurLimit:             3,
			OfflineLimit:          2,
			HeartbeatIntervalSecs: 30,
			OneDevicePerSession:   true,
			ReconnectPolicy:       "allow_once",
			RiskThreshold:         0.7,
		}
		h.db.Create(&setting)
	}

	fields := map[string]interface{}{}
	if req.FullscreenRequired != nil {
		fields["fullscreen_required"] = *req.FullscreenRequired
	}
	if req.TabSwitchLimit != nil {
		fields["tab_switch_limit"] = *req.TabSwitchLimit
	}
	if req.BlurLimit != nil {
		fields["blur_limit"] = *req.BlurLimit
	}
	if req.OfflineLimit != nil {
		fields["offline_limit"] = *req.OfflineLimit
	}
	if req.HeartbeatIntervalSecs != nil {
		fields["heartbeat_interval_secs"] = *req.HeartbeatIntervalSecs
	}
	if req.OneDevicePerSession != nil {
		fields["one_device_per_session"] = *req.OneDevicePerSession
	}
	if req.ReconnectPolicy != nil {
		fields["reconnect_policy"] = *req.ReconnectPolicy
	}
	if req.RiskThreshold != nil {
		fields["risk_threshold"] = *req.RiskThreshold
	}

	if len(fields) > 0 {
		h.db.Model(&setting).Updates(fields)
	}

	h.db.First(&setting, setting.ID)
	response.Success(c, http.StatusOK, "Security settings updated", setting)
}
