package sacertificates

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateRequest struct {
	UserID      uint    `json:"user_id" binding:"required"`
	TemplateID  *uint   `json:"template_id"`
	SourceType  string  `json:"source_type" binding:"required,oneof=olympiad mock_test"`
	SourceID    uint    `json:"source_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	FullName    string  `json:"full_name"`
	ClassName   string  `json:"class_name"`
	SubjectName string  `json:"subject_name"`
	Score       float64 `json:"score"`
	MaxScore    float64 `json:"max_score"`
	Percentage  float64 `json:"percentage"`
}

type UpdateRequest struct {
	Title   *string `json:"title"`
	FileURL *string `json:"file_url"`
	Status  *string `json:"status"`
}

type RevokeRequest struct {
	Reason string `json:"reason"`
}

type ListParams struct {
	SourceType string `form:"source_type"`
	Status     string `form:"status"`
	UserID     string `form:"user_id"`
	Search     string `form:"search"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
}

type CertificateResponse struct {
	ID                uint       `json:"id"`
	UserID            uint       `json:"user_id"`
	Username          string     `json:"username,omitempty"`
	TemplateID        *uint      `json:"template_id,omitempty"`
	SourceType        string     `json:"source_type"`
	SourceID          uint       `json:"source_id"`
	CertificateNumber string     `json:"certificate_number"`
	VerificationCode  string     `json:"verification_code"`
	Title             string     `json:"title"`
	FullName          string     `json:"full_name"`
	ClassName         string     `json:"class_name"`
	SubjectName       string     `json:"subject_name"`
	Score             float64    `json:"score"`
	MaxScore          float64    `json:"max_score"`
	Percentage        float64    `json:"percentage"`
	Status            string     `json:"status"`
	FileURL           string     `json:"file_url"`
	PDFURL            string     `json:"pdf_url"`
	IssuedAt          time.Time  `json:"issued_at"`
	RevokedAt         *time.Time `json:"revoked_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

func ToResponse(c *models.Certificate) CertificateResponse {
	resp := CertificateResponse{
		ID:                c.ID,
		UserID:            c.UserID,
		TemplateID:        c.TemplateID,
		SourceType:        string(c.SourceType),
		SourceID:          c.SourceID,
		CertificateNumber: c.CertificateNumber,
		VerificationCode:  c.VerificationCode,
		Title:             c.Title,
		FullName:          c.FullName,
		ClassName:         c.ClassName,
		SubjectName:       c.SubjectName,
		Score:             c.Score,
		MaxScore:          c.MaxScore,
		Percentage:        c.Percentage,
		Status:            c.Status,
		FileURL:           c.FileURL,
		PDFURL:            c.PDFURL,
		IssuedAt:          c.IssuedAt,
		RevokedAt:         c.RevokedAt,
		CreatedAt:         c.CreatedAt,
	}
	if c.User != nil {
		resp.Username = c.User.Username
	}
	return resp
}
