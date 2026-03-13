package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nextolympservice/go-backend/config"
)

// PanelJWTClaims — admin va superadmin uchun JWT claims
type PanelJWTClaims struct {
	StaffID  uint      `json:"staff_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"` // "admin" | "superadmin"
	Type     TokenType `json:"type"`
	jwt.RegisteredClaims
}

type PanelJWTManager struct {
	cfg *config.PanelJWTConfig
}

func NewPanelJWTManager(cfg *config.PanelJWTConfig) *PanelJWTManager {
	return &PanelJWTManager{cfg: cfg}
}

func (j *PanelJWTManager) GenerateTokenPair(staffID uint, username, role string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.generateToken(staffID, username, role, AccessToken, j.cfg.AccessSecret, j.cfg.AccessExpiryHours)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err = j.generateToken(staffID, username, role, RefreshToken, j.cfg.RefreshSecret, j.cfg.RefreshExpiryHours)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return accessToken, refreshToken, nil
}

func (j *PanelJWTManager) ValidateAccessToken(tokenString string) (*PanelJWTClaims, error) {
	return j.validateToken(tokenString, j.cfg.AccessSecret, AccessToken)
}

func (j *PanelJWTManager) ValidateRefreshToken(tokenString string) (*PanelJWTClaims, error) {
	return j.validateToken(tokenString, j.cfg.RefreshSecret, RefreshToken)
}

func (j *PanelJWTManager) generateToken(staffID uint, username, role string, tokenType TokenType, secret string, expiryHours int) (string, error) {
	claims := PanelJWTClaims{
		StaffID:  staffID,
		Username: username,
		Role:     role,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "nextolympservice-panel",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (j *PanelJWTManager) validateToken(tokenString, secret string, expectedType TokenType) (*PanelJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PanelJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(*PanelJWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	if claims.Type != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s", expectedType)
	}
	return claims, nil
}
