package adminresults

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

type ResultsHandler struct {
	db *gorm.DB
}

func NewResultsHandler(db *gorm.DB) *ResultsHandler {
	return &ResultsHandler{db: db}
}

// List GET /api/v1/admin/results
func (h *ResultsHandler) List(c *gin.Context) {
	response.Success(c, http.StatusOK, "Results", gin.H{
		"data":    []interface{}{},
		"message": "Exam engine integration pending",
	})
}

// GetByID GET /api/v1/admin/results/:id
func (h *ResultsHandler) GetByID(c *gin.Context) {
	response.Success(c, http.StatusOK, "Result", gin.H{
		"message": "Exam engine integration pending",
	})
}
