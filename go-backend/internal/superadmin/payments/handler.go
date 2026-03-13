package payments

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	repo *Repository
	db   *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
		db:   db,
	}
}

// List — to'lovlar ro'yxati
func (h *Handler) List(c *gin.Context) {
	var params ListPaymentsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.ValidationError(c, err)
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	payments, total, err := h.repo.List(params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "To'lovlarni olishda xatolik")
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, "To'lovlar", payments, params.Page, params.Limit, total)
}

// GetByID — bitta to'lov
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	payment, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "To'lov topilmadi")
		return
	}

	response.Success(c, http.StatusOK, "To'lov", payment)
}

// UpdateStatus — to'lov statusini o'zgartirish
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	payment, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "To'lov topilmadi")
		return
	}

	var req UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	oldStatus := payment.Status
	payment.Status = models.PaymentStatus(req.Status)

	if err := h.repo.Update(payment); err != nil {
		response.Error(c, http.StatusInternalServerError, "Statusni yangilashda xatolik")
		return
	}

	if req.Status == "completed" && oldStatus != models.PaymentStatusCompleted {
		h.activateRegistration(payment)
	}

	if req.Status == "refunded" && oldStatus == models.PaymentStatusCompleted {
		h.processRefund(payment)
	}

	response.Success(c, http.StatusOK, "To'lov statusi yangilandi", payment)
}

// Approve — qo'lda tasdiqlash
func (h *Handler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	payment, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "To'lov topilmadi")
		return
	}

	if payment.Status != models.PaymentStatusPending {
		response.Error(c, http.StatusBadRequest, "Faqat 'pending' statusdagi to'lovni tasdiqlash mumkin")
		return
	}

	payment.Status = models.PaymentStatusCompleted
	if err := h.repo.Update(payment); err != nil {
		response.Error(c, http.StatusInternalServerError, "To'lovni tasdiqlashda xatolik")
		return
	}

	h.activateRegistration(payment)
	response.Success(c, http.StatusOK, "To'lov tasdiqlandi", payment)
}

// Refund — qaytarish
func (h *Handler) Refund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	payment, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "To'lov topilmadi")
		return
	}

	if payment.Status != models.PaymentStatusCompleted {
		response.Error(c, http.StatusBadRequest, "Faqat 'completed' statusdagi to'lovni qaytarish mumkin")
		return
	}

	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	payment.Status = models.PaymentStatusRefunded
	if err := h.repo.Update(payment); err != nil {
		response.Error(c, http.StatusInternalServerError, "Qaytarishda xatolik")
		return
	}

	h.processRefund(payment)
	response.Success(c, http.StatusOK, "To'lov qaytarildi", payment)
}

// CreateManual — qo'lda to'lov yaratish
func (h *Handler) CreateManual(c *gin.Context) {
	var req ManualPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	payment := models.Payment{
		UserID:        req.UserID,
		SourceType:    models.PaymentSourceType(req.SourceType),
		SourceID:      req.SourceID,
		Amount:        req.Amount,
		Currency:      "UZS",
		Status:        models.PaymentStatusCompleted,
		TransactionID: fmt.Sprintf("MANUAL-%d-%d", req.UserID, time.Now().Unix()),
	}

	if err := h.repo.Create(&payment); err != nil {
		response.Error(c, http.StatusInternalServerError, "To'lov yaratishda xatolik")
		return
	}

	h.activateRegistration(&payment)
	response.Success(c, http.StatusCreated, "To'lov yaratildi", payment)
}

// Stats — to'lov statistikasi
func (h *Handler) Stats(c *gin.Context) {
	stats, err := h.repo.GetStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Statistikani olishda xatolik")
		return
	}

	response.Success(c, http.StatusOK, "To'lov statistikasi", stats)
}

func (h *Handler) activateRegistration(payment *models.Payment) {
	switch payment.SourceType {
	case "olympiad":
		h.db.Model(&models.OlympiadRegistration{}).
			Where("user_id = ? AND olympiad_id = ?", payment.UserID, payment.SourceID).
			Update("status", "registered")
	case "mock_test":
		h.db.Model(&models.MockTestRegistration{}).
			Where("user_id = ? AND mock_test_id = ?", payment.UserID, payment.SourceID).
			Update("status", "registered")
	}
}

func (h *Handler) processRefund(payment *models.Payment) {
	var balance models.Balance
	if err := h.db.Where("user_id = ?", payment.UserID).First(&balance).Error; err != nil {
		balance = models.Balance{UserID: payment.UserID, Amount: 0, Currency: "UZS"}
		h.db.Create(&balance)
	}

	balanceBefore := balance.Amount
	balance.Amount += payment.Amount
	h.db.Save(&balance)

	h.db.Create(&models.BalanceTransaction{
		UserID:          payment.UserID,
		Type:            "refund",
		Amount:          payment.Amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balance.Amount,
		Description:     fmt.Sprintf("To'lov qaytarildi #%d", payment.ID),
		SourceType:      "refund",
		SourceID:        &payment.ID,
		PerformedByType: "superadmin",
	})
}
