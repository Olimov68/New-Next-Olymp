package usercertificates

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CertificateResponse struct {
	ID                uint      `json:"id"`
	SourceType        string    `json:"source_type"`
	SourceID          uint      `json:"source_id"`
	CertificateNumber string    `json:"certificate_number"`
	Title             string    `json:"title"`
	FileURL           string    `json:"file_url"`
	IssuedAt          time.Time `json:"issued_at"`
	VerificationCode  string    `json:"verification_code"`
	CreatedAt         time.Time `json:"created_at"`
}

func ToCertificateResponse(c *models.Certificate) CertificateResponse {
	return CertificateResponse{
		ID:                c.ID,
		SourceType:        string(c.SourceType),
		SourceID:          c.SourceID,
		CertificateNumber: c.CertificateNumber,
		Title:             c.Title,
		FileURL:           c.FileURL,
		IssuedAt:          c.IssuedAt,
		VerificationCode:  c.VerificationCode,
		CreatedAt:         c.CreatedAt,
	}
}
