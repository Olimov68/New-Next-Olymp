package sanews

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
	q := r.db.Model(&models.Content{})
	if params.Type != "" {
		q = q.Where("type = ?", params.Type)
	}
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if params.Search != "" {
		q = q.Where("title ILIKE ? OR body ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	q.Count(&total)
	offset := (params.Page - 1) * params.PageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&list).Error
	return list, total, err
}

func (r *Repository) GetByID(id uint) (*models.Content, error) {
	var c models.Content
	err := r.db.First(&c, id).Error
	return &c, err
}

func (r *Repository) SlugExists(slug string) bool {
	var count int64
	r.db.Model(&models.Content{}).Where("slug = ?", slug).Count(&count)
	return count > 0
}

func (r *Repository) Create(c *models.Content) error {
	return r.db.Create(c).Error
}

func (r *Repository) Update(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Content{}).Where("id = ?", id).Updates(fields).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&models.Content{}, id).Error
}
