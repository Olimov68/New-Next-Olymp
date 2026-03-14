package user

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

func (r *Repository) GetProfileByUserID(userID uint) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) CreateProfile(profile *models.Profile) error {
	return r.db.Create(profile).Error
}

func (r *Repository) UpdateProfile(profile *models.Profile) error {
	return r.db.Save(profile).Error
}

func (r *Repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) ProfileExists(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Profile{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}
