package models

import "time"

// SecuritySetting — imtihon xavfsizlik sozlamalari
type SecuritySetting struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	FullscreenRequired     bool      `gorm:"default:true;not null" json:"fullscreen_required"`
	TabSwitchLimit         int       `gorm:"default:3;not null" json:"tab_switch_limit"`
	BlurLimit              int       `gorm:"default:3;not null" json:"blur_limit"`
	OfflineLimit           int       `gorm:"default:2;not null" json:"offline_limit"`
	HeartbeatIntervalSecs  int       `gorm:"default:30;not null" json:"heartbeat_interval_seconds"`
	OneDevicePerSession    bool      `gorm:"default:true;not null" json:"one_device_per_session"`
	ReconnectPolicy        string    `gorm:"size:50;default:allow_once" json:"reconnect_policy"` // allow_once, allow_always, deny
	RiskThreshold          float64   `gorm:"default:0.7;not null" json:"risk_threshold"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SecuritySetting) TableName() string { return "security_setting" }
