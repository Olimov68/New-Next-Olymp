package models

import "time"

type StaffRole string
type StaffStatus string

const (
	StaffRoleAdmin      StaffRole = "admin"
	StaffRoleSuperAdmin StaffRole = "superadmin"
)

const (
	StaffStatusActive  StaffStatus = "active"
	StaffStatusBlocked StaffStatus = "blocked"
)

type StaffUser struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	Username     string      `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash string      `gorm:"size:255;not null" json:"-"`
	FullName     string      `gorm:"size:200;not null" json:"full_name"`
	FirstName    string      `gorm:"size:100" json:"first_name"`
	LastName     string      `gorm:"size:100" json:"last_name"`
	Email        string      `gorm:"size:200" json:"email"`
	Phone        string      `gorm:"size:20" json:"phone"`
	Role         StaffRole   `gorm:"size:20;not null" json:"role"`
	Status       StaffStatus `gorm:"size:20;default:active;not null" json:"status"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
