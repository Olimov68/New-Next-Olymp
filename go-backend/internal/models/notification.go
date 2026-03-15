package models

import "time"

// Notification — bildirishnomalar
type Notification struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	Type       string     `gorm:"size:50;not null;index" json:"type"` // payment_success | payment_failed | registration_confirmed | result_published | certificate_ready | system | welcome
	Title      string     `gorm:"size:300;not null" json:"title"`
	Message    string     `gorm:"type:text;not null" json:"message"`
	IsRead     bool       `gorm:"default:false;not null;index" json:"is_read"`
	ReadAt     *time.Time `json:"read_at"`
	ActionURL  string     `gorm:"size:500" json:"action_url"`
	SourceType string     `gorm:"size:50" json:"source_type"` // olympiad | mock_test | payment | certificate | system
	SourceID   *uint      `json:"source_id"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Notification) TableName() string { return "notification" }
