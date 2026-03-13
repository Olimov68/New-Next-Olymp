package userfeedback

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

// Create — feedback yuborish
// POST /api/v1/user/feedback
func (h *Handler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := h.service.Create(userID.(uint), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Feedback submitted", result)
}

// List — mening feedbacklarim
// GET /api/v1/user/feedback
func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("userID")

	result, err := h.service.GetMyFeedbacks(userID.(uint))
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, http.StatusOK, "My feedbacks", result)
}

// GetByID — bitta feedback
// GET /api/v1/user/feedback/:id
func (h *Handler) GetByID(c *gin.Context) {
	userID, _ := c.Get("userID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	result, err := h.service.GetByID(uint(id), userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Feedback detail", result)
}
