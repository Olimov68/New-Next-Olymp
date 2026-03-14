package sacertificates

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateRequest struct {
	UserID            uint   `json:"user_id" binding:"required"`
	SourceType        string `json:"source_type" binding:"required,oneof=olympiad mock_test"`
	SourceID          uint   `json:"source_id" binding:"required"`
	CertificateNumber string `json:"certificate_number" binding:"required"`
	Title             string `json:"title" binding:"required"`
	FileURL           string `json:"file_url"`
}

type UpdateRequest struct {
	Title   *string `json:"title"`
	FileURL *string `json:"file_url"`
}

type ListParams struct {
	SourceType string `form:"source_type"`
	UserID     string `form:"user_id"`
	Search     string `form:"search"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
}

type CertificateResponse struct {
	ID                uint      `json:"id"`
	UserID            uint      `json:"user_id"`
	Username          string    `json:"username,omitempty"`
	SourceType        string    `json:"source_type"`
	SourceID          uint      `json:"source_id"`
	CertificateNumber string    `json:"certificate_number"`
	Title             string    `json:"title"`
	FileURL           string    `json:"file_url"`
	IssuedAt          time.Time `json:"issued_at"`
	VerificationCode  string    `json:"verification_code"`
	CreatedAt         time.Time `json:"created_at"`
}

func ToResponse(c *models.Certificate) CertificateResponse {
	resp := CertificateResponse{
		ID:                c.ID,
		UserID:            c.UserID,
		SourceType:        string(c.SourceType),
		SourceID:          c.SourceID,
		CertificateNumber: c.CertificateNumber,
		Title:             c.Title,
		FileURL:           c.FileURL,
		IssuedAt:          c.IssuedAt,
		VerificationCode:  c.VerificationCode,
		CreatedAt:         c.CreatedAt,
	}
	if c.User != nil {
		resp.Username = c.User.Username
	}
	return resp
}
