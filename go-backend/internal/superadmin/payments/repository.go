package payments

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

func (r *Repository) GetByID(id uint) (*models.Payment, error) {
	var p models.Payment
	err := r.db.Preload("User").First(&p, id).Error
	return &p, err
}

func (r *Repository) List(params ListPaymentsParams) ([]models.Payment, int64, error) {
	query := r.db.Model(&models.Payment{})

	if params.UserID > 0 {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.SourceType != "" {
		query = query.Where("source_type = ?", params.SourceType)
	}
	if params.SourceID > 0 {
		query = query.Where("source_id = ?", params.SourceID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.DateFrom != "" {
		query = query.Where("created_at >= ?", params.DateFrom)
	}
	if params.DateTo != "" {
		query = query.Where("created_at <= ?", params.DateTo+" 23:59:59")
	}
	if params.Search != "" {
		query = query.Where("transaction_id ILIKE ?", "%"+params.Search+"%")
	}

	var total int64
	query.Count(&total)

	var payments []models.Payment
	offset := (params.Page - 1) * params.Limit
	err := query.Preload("User").
		Order(params.SortBy + " " + params.SortOrder).
		Offset(offset).Limit(params.Limit).
		Find(&payments).Error

	return payments, total, err
}

func (r *Repository) Update(p *models.Payment) error {
	return r.db.Save(p).Error
}

func (r *Repository) Create(p *models.Payment) error {
	return r.db.Create(p).Error
}

func (r *Repository) GetStats() (map[string]interface{}, error) {
	var totalRevenue float64
	r.db.Model(&models.Payment{}).Where("status = ?", "completed").Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	var totalCount int64
	r.db.Model(&models.Payment{}).Count(&totalCount)

	var pendingCount int64
	r.db.Model(&models.Payment{}).Where("status = ?", "pending").Count(&pendingCount)

	var completedCount int64
	r.db.Model(&models.Payment{}).Where("status = ?", "completed").Count(&completedCount)

	var failedCount int64
	r.db.Model(&models.Payment{}).Where("status = ?", "failed").Count(&failedCount)

	var refundedCount int64
	r.db.Model(&models.Payment{}).Where("status = ?", "refunded").Count(&refundedCount)

	var refundedAmount float64
	r.db.Model(&models.Payment{}).Where("status = ?", "refunded").Select("COALESCE(SUM(amount), 0)").Scan(&refundedAmount)

	return map[string]interface{}{
		"total_revenue":   totalRevenue,
		"total_count":     totalCount,
		"pending_count":   pendingCount,
		"completed_count": completedCount,
		"failed_count":    failedCount,
		"refunded_count":  refundedCount,
		"refunded_amount": refundedAmount,
	}, nil
}
