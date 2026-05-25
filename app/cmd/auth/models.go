package auth

import "github.com/golang-jwt/jwt/v5"

const (
	TypeAccessToken  = "access"
	TypeRefreshToken = "refresh"
)

type UserClaims struct {
	UserId     uint   `json:"user_id"`
	BusinessId *uint  `json:"business_id"`
	TokenType  string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
