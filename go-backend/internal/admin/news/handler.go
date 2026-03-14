package adminnews

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	repo *Repository
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{repo: NewRepository(db)}
}

// List GET /api/v1/admin/news
func (h *Handler) List(c *gin.Context) {
	var params ListParams
	c.ShouldBindQuery(&params)
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	list, total, err := h.repo.List(params)
	if err != nil {
		response.InternalError(c)
		return
	}

	items := make([]ContentResponse, len(list))
	for i, c := range list {
		items[i] = ToContentResponse(&c)
	}
	response.Success(c, http.StatusOK, "Contents", gin.H{
		"data": items, "total": total,
		"page": params.Page, "page_size": params.PageSize,
	})
}

// GetByID GET /api/v1/admin/news/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	content, err := h.repo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "Content not found")
		return
	}
	response.Success(c, http.StatusOK, "Content", ToContentResponse(content))
}

// Create POST /api/v1/admin/news
func (h *Handler) Create(c *gin.Context) {
	staffID, _ := c.Get("staffID")
	var req CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	status := models.ContentStatusDraft
	if req.Status != "" {
		status = models.ContentStatus(req.Status)
	}

	sid := staffID.(uint)
	content := &models.Content{
		Title:       req.Title,
		Slug:        strings.ToLower(strings.ReplaceAll(req.Title, " ", "-")),
		Body:        req.Body,
		Excerpt:     req.Excerpt,
		CoverImage:  req.CoverImage,
		Type:        models.ContentType(req.Type),
		Status:      status,
		CreatedByID: &sid,
	}

	if status == models.ContentStatusPublished {
		t := PublishedAt()
		content.PublishedAt = t
	}

	if err := h.repo.Create(content); err != nil {
		response.InternalError(c)
		return
	}

	slug := fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(req.Title, " ", "-")), content.ID)
	h.repo.Update(content.ID, map[string]interface{}{"slug": slug})
	content.Slug = slug

	response.Success(c, http.StatusCreated, "Content created", ToContentResponse(content))
}

// Update PUT /api/v1/admin/news/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	existing, err := h.repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Content not found")
			return
		}
		response.InternalError(c)
		return
	}
	_ = existing

	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	fields := map[string]interface{}{}
	if req.Title != nil {
		fields["title"] = *req.Title
		fields["slug"] = fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(*req.Title, " ", "-")), id)
	}
	if req.Body != nil {
		fields["body"] = *req.Body
	}
	if req.Excerpt != nil {
		fields["excerpt"] = *req.Excerpt
	}
	if req.CoverImage != nil {
		fields["cover_image"] = *req.CoverImage
	}
	if req.Type != nil {
		fields["type"] = *req.Type
	}
	if req.Status != nil {
		fields["status"] = *req.Status
		if *req.Status == string(models.ContentStatusPublished) {
			fields["published_at"] = PublishedAt()
		}
	}

	if err := h.repo.Update(uint(id), fields); err != nil {
		response.InternalError(c)
		return
	}

	updated, _ := h.repo.GetByID(uint(id))
	response.Success(c, http.StatusOK, "Content updated", ToContentResponse(updated))
}

// Delete DELETE /api/v1/admin/news/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	if err := h.repo.Delete(uint(id)); err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, http.StatusOK, "Content deleted", nil)
}
