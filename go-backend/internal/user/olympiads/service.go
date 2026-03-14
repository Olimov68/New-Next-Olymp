package userolympiads

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

func (s *Service) List(params ListParams) (*PaginatedOlympiads, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	list, total, err := s.repo.List(params)
	if err != nil {
		return nil, err
	}

	items := make([]OlympiadResponse, len(list))
	for i, o := range list {
		items[i] = ToOlympiadResponse(&o)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize != 0 {
		totalPages++
	}

	return &PaginatedOlympiads{
		Data:       items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetByID(id uint) (*OlympiadResponse, error) {
	o, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("olympiad not found")
		}
		return nil, err
	}
	res := ToOlympiadResponse(o)
	return &res, nil
}

func (s *Service) GetMyOlympiads(userID uint) ([]map[string]interface{}, error) {
	regs, err := s.repo.GetMyOlympiads(userID)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(regs))
	for i, reg := range regs {
		item := map[string]interface{}{
			"id":          reg.ID,
			"olympiad_id": reg.OlympiadID,
			"status":      string(reg.Status),
			"joined_at":   reg.JoinedAt,
		}
		if reg.Olympiad != nil {
			item["olympiad"] = ToOlympiadResponse(reg.Olympiad)
		}
		result[i] = item
	}
	return result, nil
}

func (s *Service) Join(userID, olympiadID uint) (*RegistrationResponse, error) {
	olympiad, err := s.repo.GetByID(olympiadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("olympiad not found")
		}
		return nil, err
	}

	if olympiad.Status != models.OlympiadStatusPublished && olympiad.Status != models.OlympiadStatusActive {
		return nil, errors.New("this olympiad is not accepting registrations")
	}

	// Oldin qo'shilgan-qo'shilmaganligini tekshirish
	_, err = s.repo.GetRegistration(userID, olympiadID)
	if err == nil {
		return nil, errors.New("you have already joined this olympiad")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	reg := &models.OlympiadRegistration{
		UserID:     userID,
		OlympiadID: olympiadID,
		Status:     models.OlympiadRegStatusRegistered,
	}

	if err := s.repo.CreateRegistration(reg); err != nil {
		return nil, errors.New("failed to join olympiad")
	}

	return &RegistrationResponse{
		ID:         reg.ID,
		OlympiadID: reg.OlympiadID,
		Status:     string(reg.Status),
		JoinedAt:   reg.JoinedAt,
	}, nil
}
