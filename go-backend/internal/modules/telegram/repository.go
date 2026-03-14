package telegram

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCode(code *models.TelegramCode) error {
	return r.db.Create(code).Error
}

func (r *Repository) GetValidCode(code string) (*models.TelegramCode, error) {
	var tc models.TelegramCode
	err := r.db.Where("code = ? AND used = false AND expires_at > ?", code, time.Now()).First(&tc).Error
	return &tc, err
}

func (r *Repository) MarkCodeUsed(id uint) error {
	return r.db.Model(&models.TelegramCode{}).Where("id = ?", id).Update("used", true).Error
}

// DeleteTelegramCodes removes all unused codes for a telegram user (cleanup before generating new code)
func (r *Repository) DeleteTelegramCodes(telegramID int64) error {
	return r.db.Where("telegram_id = ? AND used = false", telegramID).Delete(&models.TelegramCode{}).Error
}

func (r *Repository) CreateTelegramLink(link *models.TelegramLink) error {
	return r.db.Create(link).Error
}

func (r *Repository) GetTelegramLinkByUserID(userID uint) (*models.TelegramLink, error) {
	var link models.TelegramLink
	err := r.db.Where("user_id = ?", userID).First(&link).Error
	return &link, err
}

func (r *Repository) GetTelegramLinkByTelegramID(telegramID int64) (*models.TelegramLink, error) {
	var link models.TelegramLink
	err := r.db.Where("telegram_id = ?", telegramID).First(&link).Error
	return &link, err
}

func (r *Repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *Repository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
