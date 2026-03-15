package models

import (
	"time"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
	UserStatusDeleted UserStatus = "deleted"
)

type User struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Username           string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash       string     `gorm:"size:255;not null" json:"-"`
	Status             UserStatus `gorm:"size:20;default:active;not null" json:"status"`
	IsProfileCompleted bool       `gorm:"default:false;not null" json:"is_profile_completed"`
	IsTelegramLinked   bool       `gorm:"default:false;not null" json:"is_telegram_linked"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Profile      *Profile      `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	TelegramLink *TelegramLink `gorm:"foreignKey:UserID" json:"telegram_link,omitempty"`
}

func (User) TableName() string { return "user" }
