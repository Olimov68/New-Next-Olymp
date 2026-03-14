package usercertificates

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

// List — mening sertifikatlarim
// GET /api/v1/user/certificates
func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("userID")

	result, err := h.service.GetMyCertificates(userID.(uint))
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, http.StatusOK, "My certificates", result)
}

// GetByID — bitta sertifikat
// GET /api/v1/user/certificates/:id
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

	response.Success(c, http.StatusOK, "Certificate detail", result)
}
