package models

import "time"

// AuditLog — tizim audit logi skeleton
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ActorID    uint      `gorm:"index" json:"actor_id"`
	ActorType  string    `gorm:"size:30" json:"actor_type"` // "user", "admin", "superadmin"
	Action     string    `gorm:"size:100;not null" json:"action"`
	Resource   string    `gorm:"size:100" json:"resource"`
	ResourceID *uint     `json:"resource_id,omitempty"`
	IPAddress  string    `gorm:"size:50" json:"ip_address"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	Details    string    `gorm:"type:text" json:"details"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
