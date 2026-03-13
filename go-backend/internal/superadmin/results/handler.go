package saresults

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
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// List — barcha natijalar (olympiad + mock test)
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sourceType := c.Query("type")   // olympiad | mock_test
	subject := c.Query("subject")
	userID := c.Query("user_id")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }

	type ResultRow struct {
		ID         uint    `json:"id"`
		UserID     uint    `json:"user_id"`
		Username   string  `json:"username"`
		FullName   string  `json:"full_name"`
		Type       string  `json:"type"`
		SourceID   uint    `json:"source_id"`
		SourceName string  `json:"source_name"`
		Subject    string  `json:"subject"`
		Score      float64 `json:"score"`
		MaxScore   float64 `json:"max_score"`
		Percentage float64 `json:"percentage"`
		Correct    int     `json:"correct"`
		Wrong      int     `json:"wrong"`
		Unanswered int     `json:"unanswered"`
		TimeTaken  int     `json:"time_taken"`
		Rank       int     `json:"rank"`
		Status     string  `json:"status"`
		Date       string  `json:"date"`
	}

	var allResults []ResultRow

	// Mock test natijalari
	if sourceType == "" || sourceType == "mock_test" {
		query := h.db.Model(&models.MockAttempt{}).
			Joins("LEFT JOIN users ON users.id = mock_attempts.user_id").
			Joins("LEFT JOIN profiles ON profiles.user_id = users.id").
			Joins("LEFT JOIN mock_tests ON mock_tests.id = mock_attempts.mock_test_id").
			Where("mock_attempts.status IN ?", []string{"completed", "timed_out"})

		if subject != "" {
			query = query.Where("mock_tests.subject = ?", subject)
		}
		if userID != "" {
			query = query.Where("mock_attempts.user_id = ?", userID)
		}
		if search != "" {
			query = query.Where("users.username ILIKE ? OR profiles.first_name ILIKE ? OR mock_tests.title ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		var mockResults []struct {
			models.MockAttempt
			Username  string `gorm:"column:username"`
			FirstName string `gorm:"column:first_name"`
			LastName  string `gorm:"column:last_name"`
			Title     string `gorm:"column:title"`
			Subject   string `gorm:"column:subject"`
		}
		query.Select("mock_attempts.*, users.username, profiles.first_name, profiles.last_name, mock_tests.title, mock_tests.subject").
			Find(&mockResults)

		for _, r := range mockResults {
			fullName := ""
			if r.FirstName != "" {
				fullName = r.FirstName + " " + r.LastName
			}
			allResults = append(allResults, ResultRow{
				ID: r.ID, UserID: r.UserID, Username: r.Username, FullName: fullName,
				Type: "mock_test", SourceID: r.MockTestID, SourceName: r.Title, Subject: r.Subject,
				Score: r.Score, MaxScore: r.MaxScore, Percentage: r.Percentage,
				Correct: r.Correct, Wrong: r.Wrong, Unanswered: r.Unanswered,
				TimeTaken: r.TimeTaken, Status: r.Status,
				Date: r.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
	}

	// Olympiad natijalari
	if sourceType == "" || sourceType == "olympiad" {
		query := h.db.Model(&models.OlympiadAttempt{}).
			Joins("LEFT JOIN users ON users.id = olympiad_attempts.user_id").
			Joins("LEFT JOIN profiles ON profiles.user_id = users.id").
			Joins("LEFT JOIN olympiads ON olympiads.id = olympiad_attempts.olympiad_id").
			Where("olympiad_attempts.status IN ?", []string{"completed", "timed_out"})

		if subject != "" {
			query = query.Where("olympiads.subject = ?", subject)
		}
		if userID != "" {
			query = query.Where("olympiad_attempts.user_id = ?", userID)
		}
		if search != "" {
			query = query.Where("users.username ILIKE ? OR profiles.first_name ILIKE ? OR olympiads.title ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		var olympResults []struct {
			models.OlympiadAttempt
			Username  string `gorm:"column:username"`
			FirstName string `gorm:"column:first_name"`
			LastName  string `gorm:"column:last_name"`
			Title     string `gorm:"column:title"`
			Subject   string `gorm:"column:subject"`
		}
		query.Select("olympiad_attempts.*, users.username, profiles.first_name, profiles.last_name, olympiads.title, olympiads.subject").
			Find(&olympResults)

		for _, r := range olympResults {
			fullName := ""
			if r.FirstName != "" {
				fullName = r.FirstName + " " + r.LastName
			}
			allResults = append(allResults, ResultRow{
				ID: r.ID, UserID: r.UserID, Username: r.Username, FullName: fullName,
				Type: "olympiad", SourceID: r.OlympiadID, SourceName: r.Title, Subject: r.Subject,
				Score: r.Score, MaxScore: r.MaxScore, Percentage: r.Percentage,
				Correct: r.Correct, Wrong: r.Wrong, Unanswered: r.Unanswered,
				TimeTaken: r.TimeTaken, Rank: r.Rank, Status: r.Status,
				Date: r.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
	}

	// Sort va paginate
	total := int64(len(allResults))
	_ = sortBy
	_ = sortOrder
	_ = fmt.Sprintf("") // suppress unused

	offset := (page - 1) * limit
	end := offset + limit
	if end > int(total) { end = int(total) }
	if offset > int(total) { offset = int(total) }

	response.SuccessWithPagination(c, http.StatusOK, "Natijalar", allResults[offset:end], page, limit, total)
}

// GetByID — bitta natija detali
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	resultType := c.Query("type") // olympiad | mock_test

	if resultType == "olympiad" {
		var attempt models.OlympiadAttempt
		if err := h.db.Preload("User").Preload("Olympiad").Preload("Answers.Question.Options").
			First(&attempt, id).Error; err != nil {
			response.Error(c, http.StatusNotFound, "Natija topilmadi")
			return
		}
		response.Success(c, http.StatusOK, "Olympiad natijasi", attempt)
		return
	}

	// Default: mock_test
	var attempt models.MockAttempt
	if err := h.db.Preload("User").Preload("MockTest").Preload("Answers.Question.Options").
		First(&attempt, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Natija topilmadi")
		return
	}
	response.Success(c, http.StatusOK, "Mock test natijasi", attempt)
}

// GetOlympiadRanking — olimpiada bo'yicha reyting
func (h *Handler) GetOlympiadRanking(c *gin.Context) {
	olympiadID, err := strconv.ParseUint(c.Param("olympiad_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	var attempts []models.OlympiadAttempt
	h.db.Preload("User").
		Where("olympiad_id = ? AND status IN ?", olympiadID, []string{"completed", "timed_out"}).
		Order("score DESC, time_taken ASC").
		Find(&attempts)

	// Rank berish
	type RankItem struct {
		Rank       int     `json:"rank"`
		UserID     uint    `json:"user_id"`
		Username   string  `json:"username"`
		Score      float64 `json:"score"`
		MaxScore   float64 `json:"max_score"`
		Percentage float64 `json:"percentage"`
		TimeTaken  int     `json:"time_taken"`
	}

	var ranking []RankItem
	for i, a := range attempts {
		username := ""
		if a.User != nil {
			username = a.User.Username
		}
		ranking = append(ranking, RankItem{
			Rank: i + 1, UserID: a.UserID, Username: username,
			Score: a.Score, MaxScore: a.MaxScore, Percentage: a.Percentage,
			TimeTaken: a.TimeTaken,
		})

		// Rank saqlash
		h.db.Model(&a).Update("rank", i+1)
	}

	response.Success(c, http.StatusOK, "Reyting", ranking)
}
