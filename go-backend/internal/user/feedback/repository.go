package userfeedback

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

func (r *Repository) Create(f *models.Feedback) error {
	return r.db.Create(f).Error
}

func (r *Repository) GetByUserID(userID uint) ([]models.Feedback, error) {
	var list []models.Feedback
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *Repository) GetByID(id, userID uint) (*models.Feedback, error) {
	var f models.Feedback
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&f).Error
	return &f, err
}
