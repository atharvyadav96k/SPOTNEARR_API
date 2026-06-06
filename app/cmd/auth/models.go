package auth

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TypeAccessToken  = "access"
	TypeRefreshToken = "refresh"
)

type UserClaims struct {
	UserId     uint            `json:"user_id"`
	BusinessId *uint           `json:"business_id"`
	TokenType  string          `json:"token_type"`
	UserRole   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
}
