package models

import "time"

type PaymentStatus string
type PaymentSourceType string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

const (
	PaymentSourceOlympiad PaymentSourceType = "olympiad"
	PaymentSourceMockTest PaymentSourceType = "mock_test"
)

// Payment — to'lov skeleton, kelajakda kengaytiriladi
type Payment struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	UserID       uint              `gorm:"not null;index" json:"user_id"`
	SourceType   PaymentSourceType `gorm:"size:30;not null" json:"source_type"`
	SourceID     uint              `gorm:"not null" json:"source_id"`
	Amount       float64           `gorm:"not null" json:"amount"`
	Currency     string            `gorm:"size:10;default:UZS" json:"currency"`
	Status       PaymentStatus     `gorm:"size:20;default:pending;not null" json:"status"`
	TransactionID string           `gorm:"size:200" json:"transaction_id"`
	CreatedAt    time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
