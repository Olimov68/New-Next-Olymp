package panelauth

import (
	"errors"
	"fmt"

	"github.com/nextolympservice/go-backend/internal/middleware"
	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	jwt  *utils.PanelJWTManager
}

func NewService(repo *Repository, jwt *utils.PanelJWTManager) *Service {
	return &Service{repo: repo, jwt: jwt}
}

func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	staff, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid credentials")
		}
		return nil, fmt.Errorf("failed to find staff: %w", err)
	}

	if staff.Status == models.StaffStatusBlocked {
		return nil, errors.New("Your account has been blocked")
	}

	if !utils.CheckPassword(req.Password, staff.PasswordHash) {
		return nil, errors.New("Invalid credentials")
	}

	accessToken, refreshToken, err := s.jwt.GenerateTokenPair(staff.ID, staff.Username, string(staff.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	permissions := middleware.GetStaffPermissionCodes(s.repo.GetDB(), staff.ID, string(staff.Role))

	return &LoginResponse{
		Staff: ToStaffResponse(staff),
		Tokens: TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		},
		Permissions: permissions,
		NextStep:    "dashboard",
	}, nil
}

func (s *Service) RefreshTokens(refreshTokenStr string) (*TokenPair, error) {
	claims, err := s.jwt.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	staff, err := s.repo.GetByID(claims.StaffID)
	if err != nil {
		return nil, errors.New("staff user not found")
	}

	if staff.Status == models.StaffStatusBlocked {
		return nil, errors.New("account is blocked")
	}

	accessToken, refreshToken, err := s.jwt.GenerateTokenPair(staff.ID, staff.Username, string(staff.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func (s *Service) GetMe(staffID uint) (*MeResponse, error) {
	staff, err := s.repo.GetByID(staffID)
	if err != nil {
		return nil, errors.New("staff user not found")
	}

	permissions := middleware.GetStaffPermissionCodes(s.repo.GetDB(), staff.ID, string(staff.Role))

	return &MeResponse{
		Staff:       ToStaffResponse(staff),
		Permissions: permissions,
	}, nil
}
