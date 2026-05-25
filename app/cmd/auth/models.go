package auth

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserId       uint   `json:"user_id"`
	BusinessId   *uint  `json:"business_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}
