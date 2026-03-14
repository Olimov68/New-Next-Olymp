package questions

import (
	"fmt"
	"net/http"
	"strconv"

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

// Create — yangi savol yaratish
func (h *Handler) Create(c *gin.Context) {
	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	// Kamida bitta to'g'ri javob tekshiruvi
	hasCorrect := false
	for _, o := range req.Options {
		if o.IsCorrect {
			hasCorrect = true
			break
		}
	}
	if !hasCorrect {
		response.Error(c, http.StatusBadRequest, "Kamida bitta to'g'ri javob bo'lishi shart")
		return
	}

	// Source mavjudligini tekshirish
	if err := h.validateSource(req.SourceType, req.SourceID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	question := models.Question{
		SourceType: req.SourceType,
		SourceID:   req.SourceID,
		Text:       req.Text,
		ImageURL:   req.ImageURL,
		Difficulty: req.Difficulty,
		Points:     req.Points,
		OrderNum:   req.OrderNum,
		IsActive:   true,
	}

	for _, o := range req.Options {
		question.Options = append(question.Options, models.QuestionOption{
			Label:     o.Label,
			Text:      o.Text,
			ImageURL:  o.ImageURL,
			IsCorrect: o.IsCorrect,
			OrderNum:  o.OrderNum,
		})
	}

	if err := h.repo.Create(&question); err != nil {
		response.Error(c, http.StatusInternalServerError, "Savol yaratishda xatolik")
		return
	}

	created, _ := h.repo.GetByID(question.ID)
	response.Success(c, http.StatusCreated, "Savol yaratildi", created)
}

// BulkCreate — bir nechta savol yaratish
func (h *Handler) BulkCreate(c *gin.Context) {
	var req BulkCreateQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	var questions []models.Question
	for _, qr := range req.Questions {
		// Validate each question has correct answer
		hasCorrect := false
		for _, o := range qr.Options {
			if o.IsCorrect {
				hasCorrect = true
				break
			}
		}
		if !hasCorrect {
			response.Error(c, http.StatusBadRequest, fmt.Sprintf("'%s' savoliga to'g'ri javob belgilanmagan", qr.Text[:50]))
			return
		}

		q := models.Question{
			SourceType: qr.SourceType,
			SourceID:   qr.SourceID,
			Text:       qr.Text,
			ImageURL:   qr.ImageURL,
			Difficulty: qr.Difficulty,
			Points:     qr.Points,
			OrderNum:   qr.OrderNum,
			IsActive:   true,
		}
		for _, o := range qr.Options {
			q.Options = append(q.Options, models.QuestionOption{
				Label:     o.Label,
				Text:      o.Text,
				ImageURL:  o.ImageURL,
				IsCorrect: o.IsCorrect,
				OrderNum:  o.OrderNum,
			})
		}
		questions = append(questions, q)
	}

	if err := h.repo.BulkCreate(questions); err != nil {
		response.Error(c, http.StatusInternalServerError, "Savollar yaratishda xatolik")
		return
	}

	response.Success(c, http.StatusCreated, fmt.Sprintf("%d ta savol yaratildi", len(questions)), nil)
}

// List — savollar ro'yxati
func (h *Handler) List(c *gin.Context) {
	var params ListQuestionsParams
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

	questions, total, err := h.repo.List(params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Savollarni olishda xatolik")
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, "Savollar", questions, params.Page, params.Limit, total)
}

// GetByID — bitta savolni olish
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	question, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Savol topilmadi")
		return
	}

	response.Success(c, http.StatusOK, "Savol", question)
}

// GetBySource — olympiad yoki mock_test bo'yicha barcha savollar
func (h *Handler) GetBySource(c *gin.Context) {
	sourceType := c.Query("source_type")
	sourceID, err := strconv.ParseUint(c.Query("source_id"), 10, 32)
	if err != nil || sourceType == "" {
		response.Error(c, http.StatusBadRequest, "source_type va source_id majburiy")
		return
	}

	questions, err := h.repo.ListBySource(sourceType, uint(sourceID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Savollarni olishda xatolik")
		return
	}

	response.Success(c, http.StatusOK, "Savollar", questions)
}

// Update — savolni yangilash
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	question, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Savol topilmadi")
		return
	}

	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if req.Text != "" {
		question.Text = req.Text
	}
	if req.ImageURL != "" {
		question.ImageURL = req.ImageURL
	}
	if req.Difficulty != "" {
		question.Difficulty = req.Difficulty
	}
	if req.Points > 0 {
		question.Points = req.Points
	}
	if req.OrderNum != nil {
		question.OrderNum = *req.OrderNum
	}
	if req.IsActive != nil {
		question.IsActive = *req.IsActive
	}

	if err := h.repo.Update(question); err != nil {
		response.Error(c, http.StatusInternalServerError, "Savolni yangilashda xatolik")
		return
	}

	// Options yangilash
	if len(req.Options) > 0 {
		hasCorrect := false
		for _, o := range req.Options {
			if o.IsCorrect {
				hasCorrect = true
				break
			}
		}
		if !hasCorrect {
			response.Error(c, http.StatusBadRequest, "Kamida bitta to'g'ri javob bo'lishi shart")
			return
		}

		// Eski optionlarni o'chirish va yangisini yaratish
		h.repo.DeleteOptions(question.ID)
		var options []models.QuestionOption
		for _, o := range req.Options {
			options = append(options, models.QuestionOption{
				QuestionID: question.ID,
				Label:      o.Label,
				Text:       o.Text,
				ImageURL:   o.ImageURL,
				IsCorrect:  o.IsCorrect,
				OrderNum:   o.OrderNum,
			})
		}
		h.repo.CreateOptions(options)
	}

	updated, _ := h.repo.GetByID(question.ID)
	response.Success(c, http.StatusOK, "Savol yangilandi", updated)
}

// Delete — savolni o'chirish
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	if _, err := h.repo.GetByID(uint(id)); err != nil {
		response.Error(c, http.StatusNotFound, "Savol topilmadi")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Savolni o'chirishda xatolik")
		return
	}

	response.Success(c, http.StatusOK, "Savol o'chirildi", nil)
}

func (h *Handler) validateSource(sourceType string, sourceID uint) error {
	switch sourceType {
	case "olympiad":
		var o models.Olympiad
		if err := h.db.First(&o, sourceID).Error; err != nil {
			return fmt.Errorf("Olimpiada topilmadi (ID: %d)", sourceID)
		}
	case "mock_test":
		var m models.MockTest
		if err := h.db.First(&m, sourceID).Error; err != nil {
			return fmt.Errorf("Mock test topilmadi (ID: %d)", sourceID)
		}
	}
	return nil
}
