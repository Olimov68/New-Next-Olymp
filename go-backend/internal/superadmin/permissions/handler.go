package permissions

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

// ListPermissions — barcha ruxsatlar
func (h *Handler) ListPermissions(c *gin.Context) {
	var perms []models.Permission
	module := c.Query("module")

	query := h.db.Model(&models.Permission{})
	if module != "" {
		query = query.Where("module = ?", module)
	}
	query.Order("module ASC, code ASC").Find(&perms)

	response.Success(c, http.StatusOK, "Ruxsatlar", perms)
}

// GetStaffPermissions — admin ruxsatlari
func (h *Handler) GetStaffPermissions(c *gin.Context) {
	staffID, err := strconv.ParseUint(c.Param("staff_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	var staffPerms []models.StaffPermission
	h.db.Preload("Permission").Where("staff_user_id = ?", staffID).Find(&staffPerms)

	var permCodes []string
	var permIDs []uint
	for _, sp := range staffPerms {
		if sp.Permission != nil {
			permCodes = append(permCodes, sp.Permission.Code)
			permIDs = append(permIDs, sp.PermissionID)
		}
	}

	response.Success(c, http.StatusOK, "Admin ruxsatlari", gin.H{
		"staff_id":       staffID,
		"permissions":    staffPerms,
		"codes":          permCodes,
		"permission_ids": permIDs,
	})
}

// AssignPermissions — adminlarga ruxsat berish
type AssignRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

func (h *Handler) AssignPermissions(c *gin.Context) {
	staffID, err := strconv.ParseUint(c.Param("staff_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	var staff models.StaffUser
	if err := h.db.First(&staff, staffID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Admin topilmadi")
		return
	}

	var req AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	// Eski ruxsatlarni o'chirish
	h.db.Where("staff_user_id = ?", staffID).Delete(&models.StaffPermission{})

	// Yangilarini qo'shish
	staffIDVal, _ := c.Get("staffID")
	grantedByID, _ := staffIDVal.(uint)
	for _, permID := range req.PermissionIDs {
		h.db.Create(&models.StaffPermission{
			StaffUserID:  uint(staffID),
			PermissionID: permID,
			GrantedByID:  &grantedByID,
		})
	}

	response.Success(c, http.StatusOK, "Ruxsatlar yangilandi", nil)
}

// SeedDefaultPermissions — standart ruxsatlarni yaratish (34 granular permissions)
func (h *Handler) SeedDefaults(c *gin.Context) {
	defaults := []models.Permission{
		// Olympiads (6)
		{Code: "olympiads.view", Name: "Olimpiadalarni ko'rish", Module: "olympiads", Description: "Olimpiadalarni ko'rish"},
		{Code: "olympiads.create", Name: "Olimpiada yaratish", Module: "olympiads", Description: "Yangi olimpiada yaratish"},
		{Code: "olympiads.update", Name: "Olimpiada tahrirlash", Module: "olympiads", Description: "Olimpiadani tahrirlash"},
		{Code: "olympiads.delete", Name: "Olimpiada o'chirish", Module: "olympiads", Description: "Olimpiadani o'chirish"},
		{Code: "olympiads.publish", Name: "Olimpiada nashr qilish", Module: "olympiads", Description: "Olimpiadani nashr qilish"},
		{Code: "olympiads.manage", Name: "Olimpiadalarni boshqarish", Module: "olympiads", Description: "Olimpiadalarni to'liq boshqarish"},

		// Mock Tests (6)
		{Code: "mock_tests.view", Name: "Mock testlarni ko'rish", Module: "mock_tests", Description: "Mock testlarni ko'rish"},
		{Code: "mock_tests.create", Name: "Mock test yaratish", Module: "mock_tests", Description: "Yangi mock test yaratish"},
		{Code: "mock_tests.update", Name: "Mock test tahrirlash", Module: "mock_tests", Description: "Mock testni tahrirlash"},
		{Code: "mock_tests.delete", Name: "Mock test o'chirish", Module: "mock_tests", Description: "Mock testni o'chirish"},
		{Code: "mock_tests.publish", Name: "Mock test nashr qilish", Module: "mock_tests", Description: "Mock testni nashr qilish"},
		{Code: "mock_tests.manage", Name: "Mock testlarni boshqarish", Module: "mock_tests", Description: "Mock testlarni to'liq boshqarish"},

		// News (6)
		{Code: "news.view", Name: "Yangiliklarni ko'rish", Module: "news", Description: "Yangiliklarni ko'rish"},
		{Code: "news.create", Name: "Yangilik yaratish", Module: "news", Description: "Yangi yangilik yaratish"},
		{Code: "news.update", Name: "Yangilik tahrirlash", Module: "news", Description: "Yangilikni tahrirlash"},
		{Code: "news.delete", Name: "Yangilik o'chirish", Module: "news", Description: "Yangilikni o'chirish"},
		{Code: "news.publish", Name: "Yangilik nashr qilish", Module: "news", Description: "Yangilikni nashr qilish"},
		{Code: "news.manage", Name: "Yangiliklarni boshqarish", Module: "news", Description: "Yangiliklarni to'liq boshqarish"},

		// Results (4)
		{Code: "results.view", Name: "Natijalarni ko'rish", Module: "results", Description: "Natijalarni ko'rish"},
		{Code: "results.update", Name: "Natijalarni tahrirlash", Module: "results", Description: "Natijalarni tahrirlash"},
		{Code: "results.export", Name: "Natijalarni eksport qilish", Module: "results", Description: "Natijalarni eksport qilish"},
		{Code: "results.manage", Name: "Natijalarni boshqarish", Module: "results", Description: "Natijalarni to'liq boshqarish"},

		// Certificates (6)
		{Code: "certificates.view", Name: "Sertifikatlarni ko'rish", Module: "certificates", Description: "Sertifikatlarni ko'rish"},
		{Code: "certificates.create", Name: "Sertifikat yaratish", Module: "certificates", Description: "Yangi sertifikat yaratish"},
		{Code: "certificates.update", Name: "Sertifikat tahrirlash", Module: "certificates", Description: "Sertifikatni tahrirlash"},
		{Code: "certificates.delete", Name: "Sertifikat o'chirish", Module: "certificates", Description: "Sertifikatni o'chirish"},
		{Code: "certificates.export", Name: "Sertifikat eksport qilish", Module: "certificates", Description: "Sertifikatni eksport qilish"},
		{Code: "certificates.manage", Name: "Sertifikatlarni boshqarish", Module: "certificates", Description: "Sertifikatlarni to'liq boshqarish"},

		// Users (6)
		{Code: "users.view", Name: "Foydalanuvchilarni ko'rish", Module: "users", Description: "Foydalanuvchilarni ko'rish"},
		{Code: "users.create", Name: "Foydalanuvchi yaratish", Module: "users", Description: "Yangi foydalanuvchi yaratish"},
		{Code: "users.update", Name: "Foydalanuvchi tahrirlash", Module: "users", Description: "Foydalanuvchini tahrirlash"},
		{Code: "users.delete", Name: "Foydalanuvchi o'chirish", Module: "users", Description: "Foydalanuvchini o'chirish"},
		{Code: "users.block", Name: "Foydalanuvchini bloklash", Module: "users", Description: "Foydalanuvchini bloklash/ochish"},
		{Code: "users.manage", Name: "Foydalanuvchilarni boshqarish", Module: "users", Description: "Foydalanuvchilarni to'liq boshqarish"},
	}

	created := 0
	for _, p := range defaults {
		var existing models.Permission
		if h.db.Where("code = ?", p.Code).First(&existing).Error != nil {
			h.db.Create(&p)
			created++
		}
	}

	response.Success(c, http.StatusOK, "Standart ruxsatlar yaratildi", gin.H{"created": created})
}
