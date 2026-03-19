package usermocktests

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

func (r *Repository) List(params ListParams) ([]models.MockTest, int64, error) {
	var list []models.MockTest
	var total int64

	q := r.db.Model(&models.MockTest{}).Where("status IN ?", []string{
		string(models.MockTestStatusPublished),
		string(models.MockTestStatusActive),
		string(models.MockTestStatusEnded),
	})

	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if params.Subject != "" {
		q = q.Where("subject = ?", params.Subject)
	}
	if params.Grade > 0 {
		q = q.Where("grade = ?", params.Grade)
	}
	if params.Language != "" {
		q = q.Where("language = ?", params.Language)
	}

	q.Count(&total)
	offset := (params.Page - 1) * params.PageSize
	err := q.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&list).Error
	return list, total, err
}

func (r *Repository) GetByID(id uint) (*models.MockTest, error) {
	var m models.MockTest
	err := r.db.Where("id = ? AND status IN ?", id, []string{
		string(models.MockTestStatusPublished),
		string(models.MockTestStatusActive),
		string(models.MockTestStatusEnded),
	}).First(&m).Error
	return &m, err
}

func (r *Repository) GetMyMockTests(userID uint) ([]models.MockTestRegistration, error) {
	var list []models.MockTestRegistration
	err := r.db.Preload("MockTest").Where("user_id = ?", userID).Order("joined_at DESC").Find(&list).Error
	return list, err
}

func (r *Repository) GetRegistration(userID, mockTestID uint) (*models.MockTestRegistration, error) {
	var reg models.MockTestRegistration
	err := r.db.Where("user_id = ? AND mock_test_id = ?", userID, mockTestID).First(&reg).Error
	return &reg, err
}

func (r *Repository) CountRegistrations(mockTestID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.MockTestRegistration{}).
		Where("mock_test_id = ? AND status != ?", mockTestID, string(models.MockTestRegStatusCancelled)).
		Count(&count).Error
	return count, err
}

func (r *Repository) CreateRegistration(reg *models.MockTestRegistration) error {
	return r.db.Create(reg).Error
}
