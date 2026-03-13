package usernews

import (
	"github.com/nextolympservice/go-backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(params ListParams) ([]models.Content, int64, error) {
	var list []models.Content
	var total int64

	q := r.db.Model(&models.Content{}).Where("status = ?", models.ContentStatusPublished)

	if params.Type != "" {
		q = q.Where("type = ?", params.Type)
	}

	q.Count(&total)
	offset := (params.Page - 1) * params.PageSize
	err := q.Order("published_at DESC").Offset(offset).Limit(params.PageSize).Find(&list).Error
	return list, total, err
}

func (r *Repository) GetByID(id uint) (*models.Content, error) {
	var c models.Content
	err := r.db.Where("id = ? AND status = ?", id, models.ContentStatusPublished).First(&c).Error
	return &c, err
}
