package models

import "time"

// Balance — foydalanuvchi balansi
type Balance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Amount    float64   `gorm:"default:0;not null" json:"amount"`
	Currency  string    `gorm:"size:10;default:UZS;not null" json:"currency"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BalanceTransaction — balans tranzaksiyalari tarixi
type BalanceTransaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	Type            string    `gorm:"size:20;not null" json:"type"`   // topup | deduct | refund
	Amount          float64   `gorm:"not null" json:"amount"`
	BalanceBefore   float64   `gorm:"not null" json:"balance_before"`
	BalanceAfter    float64   `gorm:"not null" json:"balance_after"`
	Description     string    `gorm:"size:500" json:"description"`
	SourceType      string    `gorm:"size:50" json:"source_type"` // payment | admin_manual | olympiad | mock_test | refund
	SourceID        *uint     `json:"source_id"`
	PerformedByID   *uint     `json:"performed_by_id"`
	PerformedByType string    `gorm:"size:20" json:"performed_by_type"` // user | admin | superadmin | system
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Balance) TableName() string { return "balance" }
func (BalanceTransaction) TableName() string { return "balance_transaction" }
