package usernews

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List — yangiliklar / e'lonlar ro'yxati
// GET /api/v1/user/news
// GET /api/v1/user/announcements  (type=announcement)
func (h *Handler) List(c *gin.Context) {
	var params ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := h.service.List(params)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, http.StatusOK, "Contents", result)
}

// GetByID — bitta content detail
// GET /api/v1/user/news/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	result, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Content detail", result)
}
