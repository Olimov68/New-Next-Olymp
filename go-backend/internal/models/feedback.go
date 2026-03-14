package models

import "time"

type FeedbackStatus string

const (
	FeedbackStatusOpen      FeedbackStatus = "open"
	FeedbackStatusInReview  FeedbackStatus = "in_review"
	FeedbackStatusAnswered  FeedbackStatus = "answered"
	FeedbackStatusClosed    FeedbackStatus = "closed"
)

type Feedback struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Category    string         `gorm:"size:100" json:"category"`
	Subject     string         `gorm:"size:300;not null" json:"subject"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	Status      FeedbackStatus `gorm:"size:20;default:open;not null" json:"status"`
	AdminReply  *string        `gorm:"type:text" json:"admin_reply,omitempty"`
	RepliedByID *uint          `json:"replied_by_id,omitempty"`
	RepliedAt   *time.Time     `json:"replied_at,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
