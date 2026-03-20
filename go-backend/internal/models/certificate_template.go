package models

import "time"

// CertificateTemplate — sertifikat shablonlari
type CertificateTemplate struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"size:200;not null" json:"name"`
	Type            string    `gorm:"size:30;not null" json:"type"` // olympiad_winner | olympiad_participant | mock_test
	Description     string    `gorm:"size:500" json:"description"`
	BackgroundImage string    `gorm:"size:500" json:"background_image"`
	LogoImage       string    `gorm:"size:500" json:"logo_image"`
	BodyTemplate    string    `gorm:"type:text" json:"body_template"`
	LayoutJSON      string    `gorm:"type:text" json:"layout_json"`
	PageSize        string    `gorm:"size:20;default:A4" json:"page_size"`
	Orientation     string    `gorm:"size:20;default:landscape" json:"orientation"`
	FontFamily      string    `gorm:"size:100;default:Arial" json:"font_family"`
	FontSize        int       `gorm:"default:16" json:"font_size"`
	FontColor       string    `gorm:"size:20;default:#000000" json:"font_color"`
	IsActive        bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedByID     *uint     `json:"created_by_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CertificateTemplate) TableName() string { return "certificate_template" }
