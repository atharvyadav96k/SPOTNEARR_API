package auth

import (
	"errors"
	"time"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/golang-jwt/jwt/v5"
)

func claimsGenerator(userID uint, businessId *uint, secret string, userRole models.UserRole, tokenType string, expire *jwt.NumericDate) UserClaims {
	return UserClaims{
		UserId:     userID,
		BusinessId: businessId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expire,
		},
		TokenType: tokenType,
		UserRole:  userRole,
	}
}

func GenerateAccessToken(userID uint, businessId *uint, secret string, userRole models.UserRole) (string, error) {
	claims := claimsGenerator(userID, businessId, secret, userRole, TypeAccessToken, jwt.NewNumericDate(
		time.Now().Add(5*time.Minute),
	))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uint, businessId *uint, secret string, userRole models.UserRole) (string, error) {
	claims := claimsGenerator(userID, businessId, secret, userRole, TypeRefreshToken, jwt.NewNumericDate(
		time.Now().Add(30*24*time.Hour),
	))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string, expectedTokenType string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	if claims.TokenType != expectedTokenType || claims.UserId == 0 {
		return nil, errors.New("invalid token type or unauthorized user payload")
	}
	return claims, nil
}
