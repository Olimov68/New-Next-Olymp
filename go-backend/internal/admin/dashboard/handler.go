package admindashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type DashboardStats struct {
	TotalUsers       int64 `json:"total_users"`
	TotalOlympiads   int64 `json:"total_olympiads"`
	TotalMockTests   int64 `json:"total_mock_tests"`
	TotalFeedbacks   int64 `json:"total_feedbacks"`
	OpenFeedbacks    int64 `json:"open_feedbacks"`
	TotalCertificates int64 `json:"total_certificates"`
}

// Stats — admin dashboard statistikasi
// GET /api/v1/admin/dashboard
func (h *Handler) Stats(c *gin.Context) {
	var stats DashboardStats
	h.db.Raw("SELECT COUNT(*) FROM users WHERE status != 'deleted'").Scan(&stats.TotalUsers)
	h.db.Raw("SELECT COUNT(*) FROM olympiads").Scan(&stats.TotalOlympiads)
	h.db.Raw("SELECT COUNT(*) FROM mock_tests").Scan(&stats.TotalMockTests)
	h.db.Raw("SELECT COUNT(*) FROM feedbacks").Scan(&stats.TotalFeedbacks)
	h.db.Raw("SELECT COUNT(*) FROM feedbacks WHERE status = 'open'").Scan(&stats.OpenFeedbacks)
	h.db.Raw("SELECT COUNT(*) FROM certificates").Scan(&stats.TotalCertificates)

	response.Success(c, http.StatusOK, "Dashboard stats", stats)
}
