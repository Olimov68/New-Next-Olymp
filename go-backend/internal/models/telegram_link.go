package models

import (
	"time"
)

// TelegramLink - alohida jadval chunki:
// 1. User va Telegram o'rtasidagi bog'lanish mustaqil entity
// 2. Kelajakda Telegram auth flow uchun qo'shimcha fieldlar kerak bo'lishi mumkin
// 3. Users jadvalini ortiqcha fieldlar bilan og'irlashtirmaydi
type TelegramLink struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	TelegramID       int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	TelegramUsername string    `gorm:"size:100" json:"telegram_username"`
	LinkedAt         time.Time `gorm:"not null" json:"linked_at"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (TelegramLink) TableName() string { return "telegram_link" }
