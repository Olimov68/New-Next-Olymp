package samocktests

import (
	"fmt"
	"strings"
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(params ListParams) ([]models.MockTest, int64, error) {
	return s.repo.List(params)
}

func (s *Service) GetByID(id uint) (*models.MockTest, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(req *CreateRequest, staffID uint) (*models.MockTest, error) {
	slug := generateSlug(req.Title)
	base := slug
	counter := 1
	for s.repo.SlugExists(slug) {
		slug = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}

	status := models.MockTestStatusDraft
	if req.Status != "" {
		status = models.MockTestStatus(req.Status)
	}
	lang := "uz"
	if req.Language != "" {
		lang = req.Language
	}
	scoring := "standard"
	if req.ScoringType != "" {
		scoring = req.ScoringType
	}

	m := &models.MockTest{
		Title:          req.Title,
		Slug:           slug,
		Description:    req.Description,
		Subject:        req.Subject,
		Grade:          req.Grade,
		Language:       lang,
		DurationMins:   req.DurationMins,
		TotalQuestions: req.TotalQuestions,
		ScoringType:    scoring,
		Status:         status,
		IsPaid:         req.IsPaid,
		Price:          req.Price,
		CreatedByID:    &staffID,
	}

	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) Update(id uint, req *UpdateRequest) (*models.MockTest, error) {
	fields := map[string]interface{}{}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if req.Subject != nil {
		fields["subject"] = *req.Subject
	}
	if req.Grade != nil {
		fields["grade"] = *req.Grade
	}
	if req.Language != nil {
		fields["language"] = *req.Language
	}
	if req.DurationMins != nil {
		fields["duration_mins"] = *req.DurationMins
	}
	if req.TotalQuestions != nil {
		fields["total_questions"] = *req.TotalQuestions
	}
	if req.ScoringType != nil {
		fields["scoring_type"] = *req.ScoringType
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.IsPaid != nil {
		fields["is_paid"] = *req.IsPaid
	}
	if req.Price != nil {
		fields["price"] = *req.Price
	}

	if err := s.repo.Update(id, fields); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func generateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		if r == ' ' || r == '-' || r == '_' {
			return '-'
		}
		return -1
	}, slug)
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("mock-test-%d", time.Now().Unix())
	}
	return slug
}
