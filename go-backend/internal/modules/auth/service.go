package auth

import (
	"errors"
	"fmt"

	"github.com/nextolympservice/go-backend/internal/models"
	"github.com/nextolympservice/go-backend/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	jwt  *utils.JWTManager
}

func NewService(repo *Repository, jwt *utils.JWTManager) *Service {
	return &Service{repo: repo, jwt: jwt}
}

// Register creates a new user account
func (s *Service) Register(req *RegisterRequest) (*RegisterResponse, error) {
	// Username validation
	if err := utils.ValidateUsername(req.Username); err != nil {
		return nil, fmt.Errorf("username: %w", err)
	}

	// Password validation
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("password: %w", err)
	}

	// Confirm password
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// Check username uniqueness
	exists, err := s.repo.UsernameExists(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:           req.Username,
		PasswordHash:       hash,
		Status:             models.UserStatusActive,
		IsProfileCompleted: false,
		IsTelegramLinked:   false,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens — registerdan keyin avtomatik login
	accessToken, refreshToken, err := s.jwt.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &RegisterResponse{
		User: ToUserResponse(user),
		Tokens: TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		},
		NextStep: DetermineNextStep(user), // yangi user uchun doim "complete_profile"
	}, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid credentials")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Blocked yoki deleted user login qila olmasin
	if user.Status == models.UserStatusBlocked {
		return nil, errors.New("Your account has been blocked")
	}
	if user.Status == models.UserStatusDeleted {
		return nil, errors.New("Invalid credentials")
	}

	// Password tekshirish
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("Invalid credentials")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.jwt.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &LoginResponse{
		User: ToUserResponse(user),
		Tokens: TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		},
		NextStep: DetermineNextStep(user),
	}, nil
}

// RefreshTokens generates a new token pair from a valid refresh token
func (s *Service) RefreshTokens(refreshTokenStr string) (*TokenPair, error) {
	claims, err := s.jwt.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// User hali mavjudligini tekshirish
	user, err := s.repo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status == models.UserStatusBlocked || user.Status == models.UserStatusDeleted {
		return nil, errors.New("account is not active")
	}

	accessToken, refreshToken, err := s.jwt.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

// GetMe returns current user data with profile
func (s *Service) GetMe(userID uint) (*MeResponse, error) {
	user, err := s.repo.GetByIDWithProfile(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &MeResponse{
		User:     ToUserResponse(user),
		Profile:  ToProfileResponse(user.Profile),
		NextStep: DetermineNextStep(user),
	}, nil
}
