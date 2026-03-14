package models

import (
	"time"
)

// TelegramCode - bot tomonidan generatsiya qilinadi, user saytga kiritadi
type TelegramCode struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TelegramID       int64     `gorm:"not null;index" json:"telegram_id"`
	TelegramUsername string    `gorm:"size:100" json:"telegram_username"`
	TelegramName     string    `gorm:"size:200" json:"telegram_name"`
	Code             string    `gorm:"size:10;uniqueIndex;not null" json:"code"`
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	Used             bool      `gorm:"default:false;not null" json:"used"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}
