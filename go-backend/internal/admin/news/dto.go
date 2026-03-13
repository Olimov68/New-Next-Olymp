package adminnews

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateContentRequest struct {
	Title      string `json:"title" binding:"required,min=3,max=500"`
	Body       string `json:"body" binding:"required"`
	Excerpt    string `json:"excerpt"`
	CoverImage string `json:"cover_image"`
	Type       string `json:"type" binding:"required,oneof=news announcement"`
	Status     string `json:"status"`
}

type UpdateContentRequest struct {
	Title      *string `json:"title"`
	Body       *string `json:"body"`
	Excerpt    *string `json:"excerpt"`
	CoverImage *string `json:"cover_image"`
	Type       *string `json:"type"`
	Status     *string `json:"status"`
}

type ContentResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Excerpt     string     `json:"excerpt"`
	CoverImage  string     `json:"cover_image"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToContentResponse(c *models.Content) ContentResponse {
	return ContentResponse{
		ID:          c.ID,
		Title:       c.Title,
		Slug:        c.Slug,
		Excerpt:     c.Excerpt,
		CoverImage:  c.CoverImage,
		Type:        string(c.Type),
		Status:      string(c.Status),
		PublishedAt: c.PublishedAt,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
