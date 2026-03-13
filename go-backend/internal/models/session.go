package models

import "time"

// Session — foydalanuvchi sessiylari
type Session struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	UserType     string     `gorm:"size:20;not null" json:"user_type"` // user | admin | superadmin
	TokenHash    string     `gorm:"size:500;not null" json:"token_hash"`
	IPAddress    string     `gorm:"size:50" json:"ip_address"`
	UserAgent    string     `gorm:"size:500" json:"user_agent"`
	DeviceInfo   string     `gorm:"size:300" json:"device_info"`
	IsActive     bool       `gorm:"default:true;not null;index" json:"is_active"`
	LastActiveAt time.Time  `json:"last_active_at"`
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// LoginAttempt — login urinishlari
type LoginAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:100;not null;index" json:"username"`
	IPAddress string    `gorm:"size:50;not null;index" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Success   bool      `gorm:"default:false;not null" json:"success"`
	Reason    string    `gorm:"size:200" json:"reason"` // success | invalid_password | user_blocked | user_not_found
	UserType  string    `gorm:"size:20;not null" json:"user_type"` // user | staff
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
