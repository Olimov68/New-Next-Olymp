package questions

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

func (r *Repository) Create(q *models.Question) error {
	return r.db.Create(q).Error
}

func (r *Repository) BulkCreate(questions []models.Question) error {
	return r.db.Create(&questions).Error
}

func (r *Repository) GetByID(id uint) (*models.Question, error) {
	var q models.Question
	err := r.db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).First(&q, id).Error
	return &q, err
}

func (r *Repository) List(params ListQuestionsParams) ([]models.Question, int64, error) {
	query := r.db.Model(&models.Question{})

	if params.SourceType != "" {
		query = query.Where("source_type = ?", params.SourceType)
	}
	if params.SourceID > 0 {
		query = query.Where("source_id = ?", params.SourceID)
	}
	if params.Difficulty != "" {
		query = query.Where("difficulty = ?", params.Difficulty)
	}
	if params.Search != "" {
		query = query.Where("text ILIKE ?", "%"+params.Search+"%")
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	var total int64
	query.Count(&total)

	var questions []models.Question
	offset := (params.Page - 1) * params.Limit
	err := query.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Order("order_num ASC, id ASC").Offset(offset).Limit(params.Limit).Find(&questions).Error

	return questions, total, err
}

func (r *Repository) ListBySource(sourceType string, sourceID uint) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Where("source_type = ? AND source_id = ?", sourceType, sourceID).
		Order("order_num ASC").Find(&questions).Error
	return questions, err
}

func (r *Repository) Update(q *models.Question) error {
	return r.db.Save(q).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&models.Question{}, id).Error
}

func (r *Repository) DeleteOptions(questionID uint) error {
	return r.db.Where("question_id = ?", questionID).Delete(&models.QuestionOption{}).Error
}

func (r *Repository) CreateOptions(options []models.QuestionOption) error {
	return r.db.Create(&options).Error
}

func (r *Repository) CountBySource(sourceType string, sourceID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Question{}).Where("source_type = ? AND source_id = ? AND is_active = true", sourceType, sourceID).Count(&count).Error
	return count, err
}
