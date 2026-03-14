package admincertificates

import (
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

// List GET /api/v1/admin/certificates
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var list []models.Certificate
	var total int64
	offset := (page - 1) * pageSize

	h.db.Model(&models.Certificate{}).Count(&total)
	h.db.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list)

	response.Success(c, http.StatusOK, "Certificates", gin.H{
		"data": list, "total": total, "page": page, "page_size": pageSize,
	})
}

// GetByID GET /api/v1/admin/certificates/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	var cert models.Certificate
	if err := h.db.Preload("User").First(&cert, id).Error; err != nil {
		response.NotFound(c, "Certificate not found")
		return
	}

	response.Success(c, http.StatusOK, "Certificate", cert)
}
