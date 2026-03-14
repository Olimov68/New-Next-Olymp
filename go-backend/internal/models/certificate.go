package models

import "time"

type CertificateSourceType string

const (
	CertSourceOlympiad CertificateSourceType = "olympiad"
	CertSourceMockTest CertificateSourceType = "mock_test"
)

type Certificate struct {
	ID                 uint                  `gorm:"primaryKey" json:"id"`
	UserID             uint                  `gorm:"not null;index" json:"user_id"`
	SourceType         CertificateSourceType `gorm:"size:30;not null" json:"source_type"`
	SourceID           uint                  `gorm:"not null" json:"source_id"`
	CertificateNumber  string                `gorm:"uniqueIndex;size:100;not null" json:"certificate_number"`
	Title              string                `gorm:"size:500;not null" json:"title"`
	FileURL            string                `gorm:"size:500" json:"file_url"`
	IssuedAt           time.Time             `gorm:"not null" json:"issued_at"`
	VerificationCode   string                `gorm:"uniqueIndex;size:100;not null" json:"verification_code"`
	CreatedAt          time.Time             `gorm:"autoCreateTime" json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
