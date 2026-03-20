package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/pkg/response"
)

type Handler struct {
	uploadDir string
	baseURL   string
}

func NewHandler(uploadDir, baseURL string) *Handler {
	os.MkdirAll(uploadDir, 0755)
	return &Handler{uploadDir: uploadDir, baseURL: baseURL}
}

// Upload handles image file upload
// POST /api/v1/upload
func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Fayl tanlanmagan")
		return
	}

	// Max 5MB
	if file.Size > 5*1024*1024 {
		response.Error(c, http.StatusBadRequest, "Fayl hajmi 5MB dan oshmasligi kerak")
		return
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true}
	if !allowed[ext] {
		response.Error(c, http.StatusBadRequest, "Faqat rasm fayllari ruxsat etilgan (jpg, png, gif, webp, svg)")
		return
	}

	// Category subfolder
	category := c.DefaultQuery("category", "general")
	dir := filepath.Join(h.uploadDir, category)
	if err := os.MkdirAll(dir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "Faylni saqlashda xatolik")
		return
	}

	// Generate unique filename using crypto/rand
	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		response.Error(c, http.StatusInternalServerError, "Server xatosi")
		return
	}

	// Truncate original filename for use in generated name
	origName := strings.TrimSuffix(file.Filename, ext)
	if len(origName) > 20 {
		origName = origName[:20]
	}
	origName = strings.ReplaceAll(origName, " ", "_")

	filename := fmt.Sprintf("%d_%s_%s%s", time.Now().UnixNano(), origName, hex.EncodeToString(randBytes), ext)

	dst := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, http.StatusInternalServerError, "Faylni saqlashda xatolik")
		return
	}

	url := fmt.Sprintf("%s/uploads/%s/%s", h.baseURL, category, filename)
	response.Success(c, http.StatusOK, "Fayl yuklandi", gin.H{
		"url":      url,
		"filename": filename,
		"size":     file.Size,
	})
}
