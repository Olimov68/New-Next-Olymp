package panelauth

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

func (r *Repository) GetByUsername(username string) (*models.StaffUser, error) {
	var staff models.StaffUser
	err := r.db.Where("username = ?", username).First(&staff).Error
	return &staff, err
}

func (r *Repository) GetByID(id uint) (*models.StaffUser, error) {
	var staff models.StaffUser
	err := r.db.First(&staff, id).Error
	return &staff, err
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}
