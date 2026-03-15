package models

import "time"

// Permission — ruxsatnomalar
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:100;uniqueIndex;not null" json:"code"` // manage_olympiads, manage_users, view_payments...
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	Module      string    `gorm:"size:50;not null;index" json:"module"` // olympiads, mock_tests, users, news, payments, results, certificates, feedback, settings
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// StaffPermission — admin/superadmin uchun ruxsat bog'lash
type StaffPermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StaffUserID  uint      `gorm:"not null;index" json:"staff_user_id"`
	PermissionID uint      `gorm:"not null;index" json:"permission_id"`
	GrantedByID  *uint     `json:"granted_by_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	StaffUser  *StaffUser  `gorm:"foreignKey:StaffUserID" json:"staff_user,omitempty"`
	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

func (Permission) TableName() string { return "permission" }
func (StaffPermission) TableName() string { return "staff_permission" }
