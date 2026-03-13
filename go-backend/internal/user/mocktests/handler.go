package usermocktests

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

// List — mock testlar ro'yxati
// GET /api/v1/user/mock-tests
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

	response.Success(c, http.StatusOK, "Mock tests", result)
}

// GetByID — bitta mock test detail
// GET /api/v1/user/mock-tests/:id
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

	response.Success(c, http.StatusOK, "Mock test detail", result)
}

// MyMockTests — mening mock testlarim
// GET /api/v1/user/mock-tests/my
func (h *Handler) MyMockTests(c *gin.Context) {
	userID, _ := c.Get("userID")

	result, err := h.service.GetMyMockTests(userID.(uint))
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, http.StatusOK, "My mock tests", result)
}

// Join — mock testga qo'shilish
// POST /api/v1/user/mock-tests/:id/join
func (h *Handler) Join(c *gin.Context) {
	userID, _ := c.Get("userID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	result, err := h.service.Join(userID.(uint), uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Successfully joined mock test", result)
}
