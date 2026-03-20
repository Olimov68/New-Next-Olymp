package models

import "time"

// PaymeTransaction — Payme to'lov tizimi tranzaksiyasi
type PaymeTransaction struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	PaymeID     string `gorm:"size:100;uniqueIndex;not null" json:"payme_id"` // Payme's transaction ID
	PaymentID   uint   `gorm:"not null;index" json:"payment_id"`             // Our payment ID (order_id)
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	Amount      int64  `gorm:"not null" json:"amount"`    // tiyin (1 UZS = 100 tiyin)
	State       int    `gorm:"default:1" json:"state"`    // 1=created, 2=performed, -1=cancel_before, -2=cancel_after
	Reason      *int   `json:"reason"`                    // Cancel reason code
	CreateTime  int64  `gorm:"not null" json:"create_time"`  // Payme create time (milliseconds)
	PerformTime int64  `gorm:"default:0" json:"perform_time"`
	CancelTime  int64  `gorm:"default:0" json:"cancel_time"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (PaymeTransaction) TableName() string { return "payme_transaction" }

// Payme transaction states
const (
	PaymeStateCreated         = 1
	PaymeStatePerformed       = 2
	PaymeStateCancelledBefore = -1
	PaymeStateCancelledAfter  = -2
)

// Payme error codes
const (
	PaymeRPCParseError      = -32700
	PaymeRPCInternalError   = -32400
	PaymeMethodNotFound     = -32601
	PaymeTransactionNotFound = -31003
	PaymeOrderNotFound       = -31050
	PaymeOrderAlreadyPaid    = -31051
	PaymeAmountMismatch      = -31001
	PaymeCantPerform         = -31008
	PaymeCantCancel          = -31007
	PaymeTransactionTimeout  = -31008
)

// PaymeTransactionTimeout in milliseconds (12 hours)
const PaymeTimeoutMs = 43200000
