package panelauth

import "github.com/nextolympservice/go-backend/internal/models"

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type StaffResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

type LoginResponse struct {
	Staff       StaffResponse `json:"staff"`
	Tokens      TokenPair     `json:"tokens"`
	Permissions []string      `json:"permissions"`
	NextStep    string        `json:"next_step"` // always "dashboard"
}

type MeResponse struct {
	Staff       StaffResponse `json:"staff"`
	Permissions []string      `json:"permissions"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func ToStaffResponse(s *models.StaffUser) StaffResponse {
	return StaffResponse{
		ID:        s.ID,
		Username:  s.Username,
		FullName:  s.FullName,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
		Phone:     s.Phone,
		Role:      string(s.Role),
		Status:    string(s.Status),
	}
}
