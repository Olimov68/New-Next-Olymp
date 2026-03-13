package sanews

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateRequest struct {
	Title      string `json:"title" binding:"required,min=3,max=500"`
	Body       string `json:"body"`
	Excerpt    string `json:"excerpt"`
	CoverImage string `json:"cover_image"`
	Type       string `json:"type" binding:"required,oneof=news announcement"`
	Status     string `json:"status"`
}

type UpdateRequest struct {
	Title      *string `json:"title"`
	Body       *string `json:"body"`
	Excerpt    *string `json:"excerpt"`
	CoverImage *string `json:"cover_image"`
	Type       *string `json:"type"`
	Status     *string `json:"status"`
}

type ListParams struct {
	Type     string `form:"type"`
	Status   string `form:"status"`
	Search   string `form:"search"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

type ContentResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Body        string     `json:"body"`
	Excerpt     string     `json:"excerpt"`
	CoverImage  string     `json:"cover_image"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedByID *uint      `json:"created_by_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToResponse(c *models.Content) ContentResponse {
	return ContentResponse{
		ID:          c.ID,
		Title:       c.Title,
		Slug:        c.Slug,
		Body:        c.Body,
		Excerpt:     c.Excerpt,
		CoverImage:  c.CoverImage,
		Type:        string(c.Type),
		Status:      string(c.Status),
		PublishedAt: c.PublishedAt,
		CreatedByID: c.CreatedByID,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
