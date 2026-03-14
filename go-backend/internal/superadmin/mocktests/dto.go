package samocktests

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateRequest struct {
	Title          string   `json:"title" binding:"required,min=3,max=300"`
	Description    string   `json:"description"`
	Subject        string   `json:"subject" binding:"required,min=1,max=100"`
	Grade          int      `json:"grade"`
	Language       string   `json:"language"`
	DurationMins   int      `json:"duration_minutes" binding:"required,min=1"`
	TotalQuestions int      `json:"total_questions"`
	ScoringType    string   `json:"scoring_type"`
	Status         string   `json:"status"`
	IsPaid         bool     `json:"is_paid"`
	Price          *float64 `json:"price"`
}

type UpdateRequest struct {
	Title          *string  `json:"title"`
	Description    *string  `json:"description"`
	Subject        *string  `json:"subject"`
	Grade          *int     `json:"grade"`
	Language       *string  `json:"language"`
	DurationMins   *int     `json:"duration_minutes"`
	TotalQuestions *int     `json:"total_questions"`
	ScoringType    *string  `json:"scoring_type"`
	Status         *string  `json:"status"`
	IsPaid         *bool    `json:"is_paid"`
	Price          *float64 `json:"price"`
}

type ListParams struct {
	Status   string `form:"status"`
	Subject  string `form:"subject"`
	Search   string `form:"search"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

type MockTestResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	Subject        string    `json:"subject"`
	Grade          int       `json:"grade"`
	Language       string    `json:"language"`
	DurationMins   int       `json:"duration_minutes"`
	TotalQuestions int       `json:"total_questions"`
	ScoringType    string    `json:"scoring_type"`
	Status         string    `json:"status"`
	IsPaid         bool      `json:"is_paid"`
	Price          *float64  `json:"price,omitempty"`
	CreatedByID    *uint     `json:"created_by_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToResponse(m *models.MockTest) MockTestResponse {
	return MockTestResponse{
		ID:             m.ID,
		Title:          m.Title,
		Slug:           m.Slug,
		Description:    m.Description,
		Subject:        m.Subject,
		Grade:          m.Grade,
		Language:       m.Language,
		DurationMins:   m.DurationMins,
		TotalQuestions: m.TotalQuestions,
		ScoringType:    m.ScoringType,
		Status:         string(m.Status),
		IsPaid:         m.IsPaid,
		Price:          m.Price,
		CreatedByID:    m.CreatedByID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
