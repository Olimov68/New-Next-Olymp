package superadminadmins

import (
	"time"

	"github.com/nextolympservice/go-backend/internal/models"
)

type CreateAdminRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FullName        string `json:"full_name" binding:"required,min=2,max=200"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Role            string `json:"role" binding:"required,oneof=admin superadmin"`
	PermissionIDs   []uint `json:"permission_ids"`
}

type UpdateAdminRequest struct {
	FullName      *string `json:"full_name"`
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	Email         *string `json:"email"`
	Phone         *string `json:"phone"`
	Role          *string `json:"role"`
	Status        *string `json:"status"`
	PermissionIDs *[]uint `json:"permission_ids"`
}

type AdminResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToAdminResponse(s *models.StaffUser) AdminResponse {
	return AdminResponse{
		ID:        s.ID,
		Username:  s.Username,
		FullName:  s.FullName,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
		Phone:     s.Phone,
		Role:      string(s.Role),
		Status:    string(s.Status),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

type ListParams struct {
	Role     string `form:"role"`
	Status   string `form:"status"`
	Search   string `form:"search"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}
