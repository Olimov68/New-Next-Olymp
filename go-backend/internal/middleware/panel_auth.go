package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/internal/utils"
	"github.com/nextolympservice/go-backend/pkg/response"
	"gorm.io/gorm"
)

// PanelAuthRequired — panel JWT tokenni tekshiradi (admin va superadmin)
func PanelAuthRequired(jwt *utils.PanelJWTManager, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Authorization header must be: Bearer {token}")
			c.Abort()
			return
		}

		claims, err := jwt.ValidateAccessToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Staff user mavjud va active ekanligini tekshirish
		var staff models.StaffUser
		if err := db.First(&staff, claims.StaffID).Error; err != nil {
			response.Unauthorized(c, "Staff user not found")
			c.Abort()
			return
		}

		if staff.Status == models.StaffStatusBlocked {
			response.Forbidden(c, "Your account has been blocked")
			c.Abort()
			return
		}

		c.Set("staffID", staff.ID)
		c.Set("staffUsername", staff.Username)
		c.Set("staffRole", string(staff.Role))
		c.Next()
	}
}

// AdminOnly — faqat admin YOKI superadmin o'ta oladi
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("staffRole")
		r, _ := role.(string)
		if r != string(models.StaffRoleAdmin) && r != string(models.StaffRoleSuperAdmin) {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

// SuperAdminOnly — faqat superadmin o'ta oladi
func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("staffRole")
		if role != string(models.StaffRoleSuperAdmin) {
			response.Forbidden(c, "SuperAdmin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
