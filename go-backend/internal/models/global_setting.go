package models

import "time"

// GlobalSetting — platformaning umumiy sozlamalari
type GlobalSetting struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	PlatformName        string    `gorm:"size:200;default:NextOlymp" json:"platform_name"`
	DefaultLanguage     string    `gorm:"size:10;default:uz" json:"default_language"`
	SupportEmail        string    `gorm:"size:200" json:"support_email"`
	MaintenanceMode     bool      `gorm:"default:false;not null" json:"maintenance_mode"`
	RegistrationEnabled bool      `gorm:"default:true;not null" json:"registration_enabled"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (GlobalSetting) TableName() string { return "global_setting" }
