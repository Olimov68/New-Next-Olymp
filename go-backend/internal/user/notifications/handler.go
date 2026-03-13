package notifications

import (
	"fmt"
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

// List — bildirishnomalar ro'yxati
func (h *Handler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	isRead := c.Query("is_read")

	query := h.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	if isRead == "true" {
		query = query.Where("is_read = true")
	} else if isRead == "false" {
		query = query.Where("is_read = false")
	}

	var total int64
	query.Count(&total)

	var notifications []models.Notification
	offset := (page - 1) * limit
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications)

	// O'qilmagan soni
	var unreadCount int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&unreadCount)

	response.Success(c, http.StatusOK, "Bildirishnomalar", gin.H{
		"notifications": notifications,
		"total":         total,
		"unread_count":  unreadCount,
		"page":          page,
		"limit":         limit,
	})
}

// MarkAsRead — bitta bildirishnomani o'qilgan deb belgilash
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	result := h.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": gorm.Expr("NOW()")})

	if result.RowsAffected == 0 {
		response.Error(c, http.StatusNotFound, "Bildirishnoma topilmadi")
		return
	}

	response.Success(c, http.StatusOK, "O'qilgan deb belgilandi", nil)
}

// MarkAllAsRead — hammasini o'qilgan deb belgilash
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	h.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": gorm.Expr("NOW()")})

	response.Success(c, http.StatusOK, "Barcha bildirishnomalar o'qildi", nil)
}

// UnreadCount — o'qilmagan bildirishnomalar soni
func (h *Handler) UnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var count int64
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count)

	response.Success(c, http.StatusOK, "O'qilmagan soni", gin.H{"count": count})
}

// Delete — bildirishnomani o'chirish
func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Noto'g'ri ID")
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})
	if result.RowsAffected == 0 {
		response.Error(c, http.StatusNotFound, "Bildirishnoma topilmadi")
		return
	}

	response.Success(c, http.StatusOK, "Bildirishnoma o'chirildi", nil)
}

// CreateNotification — bildirishnoma yaratish helper (boshqa modullardan chaqiriladi)
func CreateNotification(db *gorm.DB, userID uint, notifType, title, message, actionURL, sourceType string, sourceID *uint) {
	notification := models.Notification{
		UserID:     userID,
		Type:       notifType,
		Title:      title,
		Message:    message,
		ActionURL:  actionURL,
		SourceType: sourceType,
		SourceID:   sourceID,
	}
	db.Create(&notification)
}
