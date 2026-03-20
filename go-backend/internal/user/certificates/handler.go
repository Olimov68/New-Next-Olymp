package usercertificates

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/pkg/response"
)

type Handler struct {
	service   *Service
	uploadDir string
}

func NewHandler(service *Service, uploadDir string) *Handler {
	return &Handler{service: service, uploadDir: uploadDir}
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

// Download — sertifikat PDFni yuklab olish
// GET /api/v1/user/certificates/:id/download
func (h *Handler) Download(c *gin.Context) {
	userID, _ := c.Get("userID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	cert, err := h.service.repo.GetByID(uint(id), userID.(uint))
	if err != nil {
		response.NotFound(c, "Sertifikat topilmadi")
		return
	}

	if cert.Status == "revoked" {
		response.Error(c, http.StatusForbidden, "Sertifikat bekor qilingan", nil)
		return
	}

	pdfPath := cert.PDFURL
	if pdfPath == "" {
		pdfPath = cert.FileURL
	}
	if pdfPath == "" {
		response.Error(c, http.StatusNotFound, "PDF fayl topilmadi", nil)
		return
	}

	fullPath := filepath.Join(h.uploadDir, filepath.Clean(pdfPath[len("/uploads"):]))
	if _, err := os.Stat(fullPath); err != nil {
		response.Error(c, http.StatusNotFound, "PDF fayl mavjud emas", nil)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+cert.CertificateNumber+".pdf")
	c.File(fullPath)
}
