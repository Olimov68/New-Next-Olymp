package telegram

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// VerifyCode - user saytda bot yuborgan kodni kiritganida chaqiriladi
func (h *Handler) VerifyCode(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Kod kiritilmadi", nil)
		return
	}

	if err := h.service.VerifyCode(userID.(uint), req.Code); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Telegram ulandi", nil)
}

// CheckStatus - userning telegram ulangan-ulanmaganini tekshiradi
func (h *Handler) CheckStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "")
		return
	}

	linked, username := h.service.CheckStatus(userID.(uint))

	response.Success(c, http.StatusOK, "Telegram status", CheckStatusResponse{
		Linked:           linked,
		TelegramUsername: username,
	})
}

// Webhook - Telegram botdan keladigan xabarlarni qabul qiladi
// Har qanday xabar kelganda yangi bir martalik kod yaratib userga yuboradi
func (h *Handler) Webhook(c *gin.Context) {
	var update TelegramUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusOK, gin.H{"ok": true}) // Telegram uchun har doim 200 qaytariladi
		return
	}

	if update.Message == nil || update.Message.From == nil {
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	msg := update.Message

	// Bot xabarlarini e'tiborsiz qoldirish
	if msg.From.IsBot {
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	telegramName := strings.TrimSpace(msg.From.FirstName + " " + msg.From.LastName)
	_ = h.service.HandleBotMessage(msg.From.ID, msg.From.Username, telegramName)

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
