package safeedback

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

func (r *Repository) List(params ListParams) ([]models.Feedback, int64, error) {
	var list []models.Feedback
	var total int64
	q := r.db.Model(&models.Feedback{})
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if params.Category != "" {
		q = q.Where("category = ?", params.Category)
	}
	if params.Search != "" {
		q = q.Where("subject ILIKE ? OR message ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	q.Count(&total)
	offset := (params.Page - 1) * params.PageSize
	err := q.Preload("User").Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&list).Error
	return list, total, err
}

func (r *Repository) GetByID(id uint) (*models.Feedback, error) {
	var f models.Feedback
	err := r.db.Preload("User").First(&f, id).Error
	return &f, err
}

func (r *Repository) Update(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Feedback{}).Where("id = ?", id).Updates(fields).Error
}
