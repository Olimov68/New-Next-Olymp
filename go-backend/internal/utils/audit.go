package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/models"
	"gorm.io/gorm"
)

func LogAudit(db *gorm.DB, c *gin.Context, action, resource string, resourceID *uint, details string) {
	staffID, _ := c.Get("staffID")
	role, _ := c.Get("staffRole")

	sid, _ := staffID.(uint)
	roleStr, _ := role.(string)

	log := &models.AuditLog{
		ActorID:    sid,
		ActorType:  roleStr,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		Details:    details,
	}
	db.Create(log)
}
