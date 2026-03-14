package admintests

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

// --- Olympiad handlers ---

// ListOlympiads GET /api/v1/admin/tests/olympiads
func (h *Handler) ListOlympiads(c *gin.Context) {
	var params TestListParams
	c.ShouldBindQuery(&params)

	items, total, err := h.service.ListOlympiads(params)
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, http.StatusOK, "Olympiads", gin.H{
		"data": items, "total": total,
		"page": params.Page, "page_size": params.PageSize,
	})
}

// GetOlympiad GET /api/v1/admin/tests/olympiads/:id
func (h *Handler) GetOlympiad(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	result, err := h.service.GetOlympiadByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Olympiad", result)
}

// CreateOlympiad POST /api/v1/admin/tests/olympiads
func (h *Handler) CreateOlympiad(c *gin.Context) {
	staffID, _ := c.Get("staffID")
	var req CreateOlympiadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	result, err := h.service.CreateOlympiad(&req, staffID.(uint))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusCreated, "Olympiad created", result)
}

// UpdateOlympiad PUT /api/v1/admin/tests/olympiads/:id
func (h *Handler) UpdateOlympiad(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	var req UpdateOlympiadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	result, err := h.service.UpdateOlympiad(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Olympiad updated", result)
}

// DeleteOlympiad DELETE /api/v1/admin/tests/olympiads/:id
func (h *Handler) DeleteOlympiad(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	if err := h.service.DeleteOlympiad(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Olympiad deleted", nil)
}

// --- MockTest handlers ---

// ListMockTests GET /api/v1/admin/tests/mock-tests
func (h *Handler) ListMockTests(c *gin.Context) {
	var params TestListParams
	c.ShouldBindQuery(&params)

	items, total, err := h.service.ListMockTests(params)
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, http.StatusOK, "Mock tests", gin.H{
		"data": items, "total": total,
		"page": params.Page, "page_size": params.PageSize,
	})
}

// GetMockTest GET /api/v1/admin/tests/mock-tests/:id
func (h *Handler) GetMockTest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	result, err := h.service.GetMockTestByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Mock test", result)
}

// CreateMockTest POST /api/v1/admin/tests/mock-tests
func (h *Handler) CreateMockTest(c *gin.Context) {
	staffID, _ := c.Get("staffID")
	var req CreateMockTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	result, err := h.service.CreateMockTest(&req, staffID.(uint))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusCreated, "Mock test created", result)
}

// UpdateMockTest PUT /api/v1/admin/tests/mock-tests/:id
func (h *Handler) UpdateMockTest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	var req UpdateMockTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	result, err := h.service.UpdateMockTest(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Mock test updated", result)
}

// DeleteMockTest DELETE /api/v1/admin/tests/mock-tests/:id
func (h *Handler) DeleteMockTest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	if err := h.service.DeleteMockTest(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Mock test deleted", nil)
}
