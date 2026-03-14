package userfeedback

import (
	"errors"

	"github.com/nextolympservice/go-backend/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(userID uint, req *CreateFeedbackRequest) (*FeedbackResponse, error) {
	f := &models.Feedback{
		UserID:   userID,
		Category: req.Category,
		Subject:  req.Subject,
		Message:  req.Message,
		Status:   models.FeedbackStatusOpen,
	}

	if err := s.repo.Create(f); err != nil {
		return nil, errors.New("failed to submit feedback")
	}

	res := ToFeedbackResponse(f)
	return &res, nil
}

func (s *Service) GetMyFeedbacks(userID uint) ([]FeedbackResponse, error) {
	list, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]FeedbackResponse, len(list))
	for i, f := range list {
		result[i] = ToFeedbackResponse(&f)
	}
	return result, nil
}

func (s *Service) GetByID(id, userID uint) (*FeedbackResponse, error) {
	f, err := s.repo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("feedback not found")
		}
		return nil, err
	}
	res := ToFeedbackResponse(f)
	return &res, nil
}
