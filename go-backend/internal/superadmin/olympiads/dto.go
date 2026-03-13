package saolympiads

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateRequest struct {
	Title         string   `json:"title" binding:"required,min=3,max=300"`
	Description   string   `json:"description"`
	Subject       string   `json:"subject" binding:"required,min=1,max=100"`
	Grade         int      `json:"grade"`
	Language      string   `json:"language"`
	StartTime     *string  `json:"start_time"`
	EndTime       *string  `json:"end_time"`
	DurationMins  int      `json:"duration_minutes" binding:"required,min=1"`
	TotalQuestions int     `json:"total_questions"`
	Rules         string   `json:"rules"`
	Status        string   `json:"status"`
	IsPaid        bool     `json:"is_paid"`
	Price         *float64 `json:"price"`
}

type UpdateRequest struct {
	Title          *string  `json:"title"`
	Description    *string  `json:"description"`
	Subject        *string  `json:"subject"`
	Grade          *int     `json:"grade"`
	Language       *string  `json:"language"`
	StartTime      *string  `json:"start_time"`
	EndTime        *string  `json:"end_time"`
	DurationMins   *int     `json:"duration_minutes"`
	TotalQuestions *int     `json:"total_questions"`
	Rules          *string  `json:"rules"`
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

type OlympiadResponse struct {
	ID             uint       `json:"id"`
	Title          string     `json:"title"`
	Slug           string     `json:"slug"`
	Description    string     `json:"description"`
	Subject        string     `json:"subject"`
	Grade          int        `json:"grade"`
	Language       string     `json:"language"`
	StartTime      *time.Time `json:"start_time,omitempty"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	DurationMins   int        `json:"duration_minutes"`
	TotalQuestions int        `json:"total_questions"`
	Rules          string     `json:"rules"`
	Status         string     `json:"status"`
	IsPaid         bool       `json:"is_paid"`
	Price          *float64   `json:"price,omitempty"`
	CreatedByID    *uint      `json:"created_by_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func ToResponse(o *models.Olympiad) OlympiadResponse {
	return OlympiadResponse{
		ID:             o.ID,
		Title:          o.Title,
		Slug:           o.Slug,
		Description:    o.Description,
		Subject:        o.Subject,
		Grade:          o.Grade,
		Language:       o.Language,
		StartTime:      o.StartTime,
		EndTime:        o.EndTime,
		DurationMins:   o.DurationMins,
		TotalQuestions: o.TotalQuestions,
		Rules:          o.Rules,
		Status:         string(o.Status),
		IsPaid:         o.IsPaid,
		Price:          o.Price,
		CreatedByID:    o.CreatedByID,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
}
