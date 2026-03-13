package usercertificates

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

func (r *Repository) GetByUserID(userID uint) ([]models.Certificate, error) {
	var list []models.Certificate
	err := r.db.Where("user_id = ?", userID).Order("issued_at DESC").Find(&list).Error
	return list, err
}

func (r *Repository) GetByID(id, userID uint) (*models.Certificate, error) {
	var c models.Certificate
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&c).Error
	return &c, err
}
