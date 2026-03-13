package exams

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/internal/user/notifications"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// ============================================
// MOCK TEST TOPSHIRISH
// ============================================

// StartMockTest — mock testni boshlash
func (h *Handler) StartMockTest(c *gin.Context) {
	userID := c.GetUint("user_id")
	mockTestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	// Mock test mavjudligini tekshirish
	var mockTest models.MockTest
	if err := h.db.First(&mockTest, mockTestID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Mock test topilmadi")
		return
	}

	if mockTest.Status != "active" && mockTest.Status != "published" {
		response.Error(c, http.StatusBadRequest, "Bu mock test hozirda faol emas")
		return
	}

	// Ro'yxatdan o'tganligini tekshirish
	var reg models.MockTestRegistration
	if err := h.db.Where("user_id = ? AND mock_test_id = ? AND status = ?", userID, mockTestID, "registered").First(&reg).Error; err != nil {
		response.Error(c, http.StatusForbidden, "Siz bu test uchun ro'yxatdan o'tmagansiz")
		return
	}

	// Faol urinish bormi tekshirish
	var activeAttempt models.MockAttempt
	if err := h.db.Where("user_id = ? AND mock_test_id = ? AND status = ?", userID, mockTestID, "in_progress").First(&activeAttempt).Error; err == nil {
		// Faol urinish bor — davom ettirish
		var answers []models.MockAttemptAnswer
		h.db.Where("attempt_id = ?", activeAttempt.ID).Find(&answers)

		response.Success(c, http.StatusOK, "Faol urinish topildi", gin.H{
			"attempt":       activeAttempt,
			"answers_given": len(answers),
			"time_left":     h.calculateTimeLeft(activeAttempt.StartedAt, mockTest.DurationMins),
		})
		return
	}

	// Savollarni olish
	var questions []models.Question
	h.db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Where("source_type = ? AND source_id = ? AND is_active = true", "mock_test", mockTestID).
		Order("order_num ASC").Find(&questions)

	if len(questions) == 0 {
		response.Error(c, http.StatusBadRequest, "Bu testda savollar mavjud emas")
		return
	}

	// Yangi urinish yaratish
	attempt := models.MockAttempt{
		UserID:     userID,
		MockTestID: uint(mockTestID),
		StartedAt:  time.Now(),
		Status:     "in_progress",
		MaxScore:   h.calculateMaxScore(questions),
	}
	h.db.Create(&attempt)

	// Savollardan is_correct ni yashirish
	safeQuestions := h.hideCorrectAnswers(questions)

	response.Success(c, http.StatusCreated, "Test boshlandi", gin.H{
		"attempt_id":    attempt.ID,
		"questions":     safeQuestions,
		"total":         len(questions),
		"duration_mins": mockTest.DurationMins,
		"started_at":    attempt.StartedAt,
		"ends_at":       attempt.StartedAt.Add(time.Duration(mockTest.DurationMins) * time.Minute),
	})
}

// SubmitMockAnswer — bitta savolga javob berish
type SubmitAnswerRequest struct {
	QuestionID uint  `json:"question_id" binding:"required,min=1"`
	OptionID   *uint `json:"option_id"` // nil = javob berilmadi
}

func (h *Handler) SubmitMockAnswer(c *gin.Context) {
	userID := c.GetUint("user_id")
	attemptID, err := strconv.ParseUint(c.Param("attempt_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri attempt ID")
		return
	}

	var attempt models.MockAttempt
	if err := h.db.Where("id = ? AND user_id = ? AND status = ?", attemptID, userID, "in_progress").First(&attempt).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Faol urinish topilmadi")
		return
	}

	// Vaqt tugaganini tekshirish
	var mockTest models.MockTest
	h.db.First(&mockTest, attempt.MockTestID)
	if h.isTimeUp(attempt.StartedAt, mockTest.DurationMins) {
		h.finishMockAttempt(&attempt)
		response.Error(c, http.StatusBadRequest, "Vaqt tugadi, test yakunlandi")
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	// Oldingi javobni tekshirish
	var existing models.MockAttemptAnswer
	if err := h.db.Where("attempt_id = ? AND question_id = ?", attemptID, req.QuestionID).First(&existing).Error; err == nil {
		// Javob yangilash
		isCorrect := false
		if req.OptionID != nil {
			var opt models.QuestionOption
			if h.db.First(&opt, *req.OptionID).Error == nil {
				isCorrect = opt.IsCorrect
			}
		}
		existing.SelectedOptionID = req.OptionID
		existing.IsCorrect = isCorrect
		h.db.Save(&existing)
		response.Success(c, http.StatusOK, "Javob yangilandi", existing)
		return
	}

	// Yangi javob
	isCorrect := false
	if req.OptionID != nil {
		var opt models.QuestionOption
		if h.db.First(&opt, *req.OptionID).Error == nil {
			isCorrect = opt.IsCorrect
		}
	}

	answer := models.MockAttemptAnswer{
		AttemptID:        uint(attemptID),
		QuestionID:       req.QuestionID,
		SelectedOptionID: req.OptionID,
		IsCorrect:        isCorrect,
	}
	h.db.Create(&answer)

	response.Success(c, http.StatusCreated, "Javob saqlandi", answer)
}

// FinishMockTest — testni yakunlash
func (h *Handler) FinishMockTest(c *gin.Context) {
	userID := c.GetUint("user_id")
	attemptID, err := strconv.ParseUint(c.Param("attempt_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri attempt ID")
		return
	}

	var attempt models.MockAttempt
	if err := h.db.Where("id = ? AND user_id = ? AND status = ?", attemptID, userID, "in_progress").First(&attempt).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Faol urinish topilmadi")
		return
	}

	result := h.finishMockAttempt(&attempt)
	response.Success(c, http.StatusOK, "Test yakunlandi", result)
}

// GetMockAttemptResult — urinish natijasini ko'rish
func (h *Handler) GetMockAttemptResult(c *gin.Context) {
	userID := c.GetUint("user_id")
	attemptID, err := strconv.ParseUint(c.Param("attempt_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri attempt ID")
		return
	}

	var attempt models.MockAttempt
	if err := h.db.Preload("MockTest").Preload("Answers.Question.Options").
		Where("id = ? AND user_id = ?", attemptID, userID).First(&attempt).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Urinish topilmadi")
		return
	}

	if attempt.Status == "in_progress" {
		response.Error(c, http.StatusBadRequest, "Test hali yakunlanmagan")
		return
	}

	response.Success(c, http.StatusOK, "Natija", attempt)
}

// GetMyMockAttempts — mening urinishlarim
func (h *Handler) GetMyMockAttempts(c *gin.Context) {
	userID := c.GetUint("user_id")
	mockTestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	var attempts []models.MockAttempt
	h.db.Where("user_id = ? AND mock_test_id = ?", userID, mockTestID).
		Order("created_at DESC").Find(&attempts)

	response.Success(c, http.StatusOK, "Urinishlar", attempts)
}

// ============================================
// OLYMPIAD TOPSHIRISH
// ============================================

// StartOlympiad — olimpiadani boshlash
func (h *Handler) StartOlympiad(c *gin.Context) {
	userID := c.GetUint("user_id")
	olympiadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	var olympiad models.Olympiad
	if err := h.db.First(&olympiad, olympiadID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Olimpiada topilmadi")
		return
	}

	if olympiad.Status != "active" {
		response.Error(c, http.StatusBadRequest, "Bu olimpiada hozirda faol emas")
		return
	}

	// Vaqt oralig'ini tekshirish
	now := time.Now()
	if olympiad.StartTime != nil && now.Before(*olympiad.StartTime) {
		response.Error(c, http.StatusBadRequest, "Olimpiada hali boshlanmagan")
		return
	}
	if olympiad.EndTime != nil && now.After(*olympiad.EndTime) {
		response.Error(c, http.StatusBadRequest, "Olimpiada tugagan")
		return
	}

	// Ro'yxatdan o'tganligini tekshirish
	var reg models.OlympiadRegistration
	if err := h.db.Where("user_id = ? AND olympiad_id = ? AND status IN ?", userID, olympiadID, []string{"registered", "participant"}).First(&reg).Error; err != nil {
		response.Error(c, http.StatusForbidden, "Siz bu olimpiada uchun ro'yxatdan o'tmagansiz")
		return
	}

	// Faol urinish
	var activeAttempt models.OlympiadAttempt
	if err := h.db.Where("user_id = ? AND olympiad_id = ? AND status = ?", userID, olympiadID, "in_progress").First(&activeAttempt).Error; err == nil {
		response.Success(c, http.StatusOK, "Faol urinish topildi", gin.H{
			"attempt":   activeAttempt,
			"time_left": h.calculateTimeLeft(activeAttempt.StartedAt, olympiad.DurationMins),
		})
		return
	}

	// Oldingi urinish bor-yo'qligini tekshirish (olympiada faqat 1 urinish)
	var prevAttempt models.OlympiadAttempt
	if err := h.db.Where("user_id = ? AND olympiad_id = ? AND status IN ?", userID, olympiadID, []string{"completed", "timed_out"}).First(&prevAttempt).Error; err == nil {
		response.Error(c, http.StatusBadRequest, "Siz bu olimpiadada allaqachon qatnashgansiz")
		return
	}

	// Savollar
	var questions []models.Question
	h.db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_num ASC")
	}).Where("source_type = ? AND source_id = ? AND is_active = true", "olympiad", olympiadID).
		Order("order_num ASC").Find(&questions)

	if len(questions) == 0 {
		response.Error(c, http.StatusBadRequest, "Bu olimpiadada savollar mavjud emas")
		return
	}

	attempt := models.OlympiadAttempt{
		UserID:     userID,
		OlympiadID: uint(olympiadID),
		StartedAt:  time.Now(),
		Status:     "in_progress",
		MaxScore:   h.calculateMaxScore(questions),
	}
	h.db.Create(&attempt)

	// Registratsiyani yangilash
	h.db.Model(&reg).Update("status", "participant")

	safeQuestions := h.hideCorrectAnswers(questions)

	response.Success(c, http.StatusCreated, "Olimpiada boshlandi", gin.H{
		"attempt_id":    attempt.ID,
		"questions":     safeQuestions,
		"total":         len(questions),
		"duration_mins": olympiad.DurationMins,
		"started_at":    attempt.StartedAt,
		"ends_at":       attempt.StartedAt.Add(time.Duration(olympiad.DurationMins) * time.Minute),
	})
}

// SubmitOlympiadAnswer — olimpiada javobi
func (h *Handler) SubmitOlympiadAnswer(c *gin.Context) {
	userID := c.GetUint("user_id")
	attemptID, err := strconv.ParseUint(c.Param("attempt_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri attempt ID")
		return
	}

	var attempt models.OlympiadAttempt
	if err := h.db.Where("id = ? AND user_id = ? AND status = ?", attemptID, userID, "in_progress").First(&attempt).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Faol urinish topilmadi")
		return
	}

	var olympiad models.Olympiad
	h.db.First(&olympiad, attempt.OlympiadID)
	if h.isTimeUp(attempt.StartedAt, olympiad.DurationMins) {
		h.finishOlympiadAttempt(&attempt)
		response.Error(c, http.StatusBadRequest, "Vaqt tugadi, olimpiada yakunlandi")
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	var existing models.OlympiadAttemptAnswer
	if err := h.db.Where("attempt_id = ? AND question_id = ?", attemptID, req.QuestionID).First(&existing).Error; err == nil {
		isCorrect := false
		if req.OptionID != nil {
			var opt models.QuestionOption
			if h.db.First(&opt, *req.OptionID).Error == nil {
				isCorrect = opt.IsCorrect
			}
		}
		existing.SelectedOptionID = req.OptionID
		existing.IsCorrect = isCorrect
		h.db.Save(&existing)
		response.Success(c, http.StatusOK, "Javob yangilandi", existing)
		return
	}

	isCorrect := false
	if req.OptionID != nil {
		var opt models.QuestionOption
		if h.db.First(&opt, *req.OptionID).Error == nil {
			isCorrect = opt.IsCorrect
		}
	}

	answer := models.OlympiadAttemptAnswer{
		AttemptID:        uint(attemptID),
		QuestionID:       req.QuestionID,
		SelectedOptionID: req.OptionID,
		IsCorrect:        isCorrect,
	}
	h.db.Create(&answer)

	response.Success(c, http.StatusCreated, "Javob saqlandi", answer)
}

// FinishOlympiad — olimpiadani yakunlash
func (h *Handler) FinishOlympiad(c *gin.Context) {
	userID := c.GetUint("user_id")
	attemptID, err := strconv.ParseUint(c.Param("attempt_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri attempt ID")
		return
	}

	var attempt models.OlympiadAttempt
	if err := h.db.Where("id = ? AND user_id = ? AND status = ?", attemptID, userID, "in_progress").First(&attempt).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Faol urinish topilmadi")
		return
	}

	result := h.finishOlympiadAttempt(&attempt)
	response.Success(c, http.StatusOK, "Olimpiada yakunlandi", result)
}

// GetOlympiadAttemptResult — olimpiada natijasi
func (h *Handler) GetOlympiadAttemptResult(c *gin.Context) {
	userID := c.GetUint("user_id")
	attemptID, err := strconv.ParseUint(c.Param("attempt_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri attempt ID")
		return
	}

	var attempt models.OlympiadAttempt
	if err := h.db.Preload("Olympiad").Preload("Answers.Question.Options").
		Where("id = ? AND user_id = ?", attemptID, userID).First(&attempt).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Urinish topilmadi")
		return
	}

	if attempt.Status == "in_progress" {
		response.Error(c, http.StatusBadRequest, "Olimpiada hali yakunlanmagan")
		return
	}

	response.Success(c, http.StatusOK, "Natija", attempt)
}

// ============================================
// HELPER FUNKSIYALAR
// ============================================

func (h *Handler) finishMockAttempt(attempt *models.MockAttempt) map[string]interface{} {
	now := time.Now()
	attempt.FinishedAt = &now
	attempt.TimeTaken = int(now.Sub(attempt.StartedAt).Seconds())

	var answers []models.MockAttemptAnswer
	h.db.Where("attempt_id = ?", attempt.ID).Find(&answers)

	var totalQuestions int64
	h.db.Model(&models.Question{}).Where("source_type = ? AND source_id = ? AND is_active = true", "mock_test", attempt.MockTestID).Count(&totalQuestions)

	correct := 0
	wrong := 0
	for _, a := range answers {
		if a.SelectedOptionID != nil {
			if a.IsCorrect {
				correct++
			} else {
				wrong++
			}
		}
	}

	unanswered := int(totalQuestions) - len(answers)
	score := float64(correct) // har bir to'g'ri javob 1 ball
	maxScore := attempt.MaxScore
	if maxScore == 0 {
		maxScore = float64(totalQuestions)
	}
	percentage := 0.0
	if maxScore > 0 {
		percentage = math.Round(score/maxScore*100*10) / 10
	}

	attempt.Correct = correct
	attempt.Wrong = wrong
	attempt.Unanswered = unanswered
	attempt.Score = score
	attempt.Percentage = percentage

	var mockTest models.MockTest
	h.db.First(&mockTest, attempt.MockTestID)
	if h.isTimeUp(attempt.StartedAt, mockTest.DurationMins) {
		attempt.Status = "timed_out"
	} else {
		attempt.Status = "completed"
	}

	h.db.Save(attempt)

	// Registration statusini yangilash
	h.db.Model(&models.MockTestRegistration{}).
		Where("user_id = ? AND mock_test_id = ?", attempt.UserID, attempt.MockTestID).
		Update("status", "completed")

	// Bildirishnoma
	notifications.CreateNotification(h.db, attempt.UserID, "result_published",
		"Test natijasi tayyor",
		fmt.Sprintf("Siz %s testida %.1f%% natija ko'rsatdingiz", mockTest.Title, percentage),
		fmt.Sprintf("/dashboard/results?attempt=%d", attempt.ID),
		"mock_test", &attempt.MockTestID)

	return map[string]interface{}{
		"attempt_id": attempt.ID,
		"score":      score,
		"max_score":  maxScore,
		"correct":    correct,
		"wrong":      wrong,
		"unanswered": unanswered,
		"percentage": percentage,
		"time_taken": attempt.TimeTaken,
		"status":     attempt.Status,
	}
}

func (h *Handler) finishOlympiadAttempt(attempt *models.OlympiadAttempt) map[string]interface{} {
	now := time.Now()
	attempt.FinishedAt = &now
	attempt.TimeTaken = int(now.Sub(attempt.StartedAt).Seconds())

	var answers []models.OlympiadAttemptAnswer
	h.db.Where("attempt_id = ?", attempt.ID).Find(&answers)

	var totalQuestions int64
	h.db.Model(&models.Question{}).Where("source_type = ? AND source_id = ? AND is_active = true", "olympiad", attempt.OlympiadID).Count(&totalQuestions)

	correct := 0
	wrong := 0
	for _, a := range answers {
		if a.SelectedOptionID != nil {
			if a.IsCorrect {
				correct++
			} else {
				wrong++
			}
		}
	}

	unanswered := int(totalQuestions) - len(answers)
	score := float64(correct)
	maxScore := attempt.MaxScore
	if maxScore == 0 {
		maxScore = float64(totalQuestions)
	}
	percentage := 0.0
	if maxScore > 0 {
		percentage = math.Round(score/maxScore*100*10) / 10
	}

	attempt.Correct = correct
	attempt.Wrong = wrong
	attempt.Unanswered = unanswered
	attempt.Score = score
	attempt.Percentage = percentage

	var olympiad models.Olympiad
	h.db.First(&olympiad, attempt.OlympiadID)
	if h.isTimeUp(attempt.StartedAt, olympiad.DurationMins) {
		attempt.Status = "timed_out"
	} else {
		attempt.Status = "completed"
	}

	h.db.Save(attempt)

	// Registration statusini yangilash
	h.db.Model(&models.OlympiadRegistration{}).
		Where("user_id = ? AND olympiad_id = ?", attempt.UserID, attempt.OlympiadID).
		Update("status", "completed")

	// Bildirishnoma
	notifications.CreateNotification(h.db, attempt.UserID, "result_published",
		"Olimpiada natijasi tayyor",
		fmt.Sprintf("Siz %s olimpiadasida %.1f%% natija ko'rsatdingiz", olympiad.Title, percentage),
		fmt.Sprintf("/dashboard/results?attempt=%d", attempt.ID),
		"olympiad", &attempt.OlympiadID)

	return map[string]interface{}{
		"attempt_id": attempt.ID,
		"score":      score,
		"max_score":  maxScore,
		"correct":    correct,
		"wrong":      wrong,
		"unanswered": unanswered,
		"percentage": percentage,
		"time_taken": attempt.TimeTaken,
		"status":     attempt.Status,
	}
}

func (h *Handler) calculateTimeLeft(startedAt time.Time, durationMins int) int {
	elapsed := time.Since(startedAt)
	total := time.Duration(durationMins) * time.Minute
	left := total - elapsed
	if left < 0 {
		return 0
	}
	return int(left.Seconds())
}

func (h *Handler) isTimeUp(startedAt time.Time, durationMins int) bool {
	return time.Since(startedAt) > time.Duration(durationMins)*time.Minute
}

func (h *Handler) calculateMaxScore(questions []models.Question) float64 {
	total := 0.0
	for _, q := range questions {
		total += q.Points
	}
	return total
}

type safeQuestion struct {
	ID         uint         `json:"id"`
	Text       string       `json:"text"`
	ImageURL   string       `json:"image_url"`
	Difficulty string       `json:"difficulty"`
	Points     float64      `json:"points"`
	OrderNum   int          `json:"order_num"`
	Options    []safeOption `json:"options"`
}

type safeOption struct {
	ID       uint   `json:"id"`
	Label    string `json:"label"`
	Text     string `json:"text"`
	ImageURL string `json:"image_url"`
	OrderNum int    `json:"order_num"`
}

func (h *Handler) hideCorrectAnswers(questions []models.Question) []safeQuestion {
	var safe []safeQuestion
	for _, q := range questions {
		sq := safeQuestion{
			ID: q.ID, Text: q.Text, ImageURL: q.ImageURL,
			Difficulty: q.Difficulty, Points: q.Points, OrderNum: q.OrderNum,
		}
		for _, o := range q.Options {
			sq.Options = append(sq.Options, safeOption{
				ID: o.ID, Label: o.Label, Text: o.Text,
				ImageURL: o.ImageURL, OrderNum: o.OrderNum,
			})
		}
		safe = append(safe, sq)
	}
	return safe
}
