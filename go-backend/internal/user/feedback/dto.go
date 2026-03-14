package userfeedback

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateFeedbackRequest struct {
	Category string `json:"category" binding:"required,max=100"`
	Subject  string `json:"subject" binding:"required,min=3,max=300"`
	Message  string `json:"message" binding:"required,min=10"`
}

type FeedbackResponse struct {
	ID         uint       `json:"id"`
	Category   string     `json:"category"`
	Subject    string     `json:"subject"`
	Message    string     `json:"message"`
	Status     string     `json:"status"`
	AdminReply *string    `json:"admin_reply,omitempty"`
	RepliedAt  *time.Time `json:"replied_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func ToFeedbackResponse(f *models.Feedback) FeedbackResponse {
	return FeedbackResponse{
		ID:         f.ID,
		Category:   f.Category,
		Subject:    f.Subject,
		Message:    f.Message,
		Status:     string(f.Status),
		AdminReply: f.AdminReply,
		RepliedAt:  f.RepliedAt,
		CreatedAt:  f.CreatedAt,
		UpdatedAt:  f.UpdatedAt,
	}
}
