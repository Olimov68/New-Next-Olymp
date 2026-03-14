package sanews

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

func (s *Service) List(params ListParams) ([]models.Content, int64, error) {
	return s.repo.List(params)
}

func (s *Service) GetByID(id uint) (*models.Content, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(req *CreateRequest, staffID uint) (*models.Content, error) {
	slug := generateSlug(req.Title)
	base := slug
	counter := 1
	for s.repo.SlugExists(slug) {
		slug = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}

	status := models.ContentStatusDraft
	if req.Status != "" {
		status = models.ContentStatus(req.Status)
	}

	c := &models.Content{
		Title:       req.Title,
		Slug:        slug,
		Body:        req.Body,
		Excerpt:     req.Excerpt,
		CoverImage:  req.CoverImage,
		Type:        models.ContentType(req.Type),
		Status:      status,
		CreatedByID: &staffID,
	}

	// Auto-set published_at when status is published
	if status == models.ContentStatusPublished {
		now := time.Now()
		c.PublishedAt = &now
	}

	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) Update(id uint, req *UpdateRequest) (*models.Content, error) {
	fields := map[string]interface{}{}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Body != nil {
		fields["body"] = *req.Body
	}
	if req.Excerpt != nil {
		fields["excerpt"] = *req.Excerpt
	}
	if req.CoverImage != nil {
		fields["cover_image"] = *req.CoverImage
	}
	if req.Type != nil {
		fields["type"] = *req.Type
	}
	if req.Status != nil {
		fields["status"] = *req.Status
		// Auto-set published_at when status changes to published
		if models.ContentStatus(*req.Status) == models.ContentStatusPublished {
			existing, err := s.repo.GetByID(id)
			if err == nil && existing.PublishedAt == nil {
				now := time.Now()
				fields["published_at"] = now
			}
		}
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
		slug = fmt.Sprintf("content-%d", time.Now().Unix())
	}
	return slug
}
