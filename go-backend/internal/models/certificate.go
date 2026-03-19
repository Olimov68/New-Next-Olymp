package models

import "time"

type CertificateSourceType string

const (
	CertSourceOlympiad CertificateSourceType = "olympiad"
	CertSourceMockTest CertificateSourceType = "mock_test"
)

// CertificateType — sertifikat turi
const (
	CertTypeOlympiad  = "olympiad"
	CertTypeMockRasch = "mock_rasch"
)

const (
	CertStatusActive  = "active"
	CertStatusRevoked = "revoked"
)

type Certificate struct {
	ID                uint                  `gorm:"primaryKey" json:"id"`
	UserID            uint                  `gorm:"not null;index" json:"user_id"`
	TemplateID        *uint                 `gorm:"index" json:"template_id,omitempty"`
	CertificateType   string                `gorm:"size:30;default:olympiad" json:"certificate_type"` // olympiad | mock_rasch
	SourceType        CertificateSourceType `gorm:"size:30;not null" json:"source_type"`
	SourceID          uint                  `gorm:"not null" json:"source_id"`
	CertificateNumber string                `gorm:"uniqueIndex;size:100;not null" json:"certificate_number"`
	VerificationCode  string                `gorm:"uniqueIndex;size:100;not null" json:"verification_code"`
	Title             string                `gorm:"size:500;not null" json:"title"`
	FullName          string                `gorm:"size:200" json:"full_name"`
	ClassName         string                `gorm:"size:100" json:"class_name"`
	SubjectName       string                `gorm:"size:200" json:"subject_name"`
	Score             float64               `json:"score"`          // raw score
	ScaledScore       float64               `json:"scaled_score"`   // Rasch scaled score
	MaxScore          float64               `json:"max_score"`
	Percentage        float64               `json:"percentage"`
	Grade             string                `gorm:"size:10" json:"grade"`   // A+, A, B, C (mock_rasch uchun)
	Rank              *int                  `json:"rank,omitempty"`         // olimpiada o'rni
	Status            string                `gorm:"size:20;default:active;not null" json:"status"`
	FileURL           string                `gorm:"size:500" json:"file_url"`
	PDFURL            string                `gorm:"size:500" json:"pdf_url"`
	IssuedAt          time.Time             `gorm:"not null" json:"issued_at"`
	ValidUntil        *time.Time            `json:"valid_until,omitempty"` // amal qilish muddati
	RevokedAt         *time.Time            `json:"revoked_at,omitempty"`
	CreatedAt         time.Time             `gorm:"autoCreateTime" json:"created_at"`

	User     *User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template *CertificateTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}

func (Certificate) TableName() string { return "certificate" }

// CalculateGrade — Rasch T-ball (T = 50 + 10Z) asosida daraja hisoblash
// Rasman baholash mezoniga asosan:
// 70+    = A+  (maks ball)
// 65-69.9 = A   (maks ball)
// 60-64.9 = B+  (proporsional ball)
// 55-59.9 = B   (proporsional ball)
// 50-54.9 = C+  (proporsional ball)
// 46-49.9 = C   (proporsional ball)
// 46 dan past = sertifikat berilmaydi
func CalculateGrade(tScore float64) string {
	switch {
	case tScore >= 70:
		return "A+"
	case tScore >= 65:
		return "A"
	case tScore >= 60:
		return "B+"
	case tScore >= 55:
		return "B"
	case tScore >= 50:
		return "C+"
	case tScore >= 46:
		return "C"
	default:
		return "" // sertifikat berilmaydi
	}
}

// IsEligibleForCertificate — sertifikat olish huquqi bormi (T-ball >= 46)
func IsEligibleForCertificate(tScore float64) bool {
	return tScore >= 46
}

// CalculateTScore — Rasch theta dan T-ball hisoblash
// Formula: T = 50 + 10 * Z, Z = (theta - mu) / sigma
func CalculateTScore(theta, mu, sigma float64) float64 {
	if sigma == 0 {
		return 50
	}
	z := (theta - mu) / sigma
	return 50 + 10*z
}
