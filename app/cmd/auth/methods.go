package auth

import (
	"errors"
	"time"

	usermodel "github.com/Developer-Aadesh/spotnearr-database/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func claimsGenerator(userID uint, businessId *uint, secret string, userRole usermodel.UserRole, tokenType string, expire *jwt.NumericDate) UserClaims {
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

func GenerateAccessToken(userID uint, businessId *uint, secret string, userRole usermodel.UserRole) (string, time.Time, error) {
	expiry := time.Now().Add(5 * time.Minute)
	claims := claimsGenerator(userID, businessId, secret, userRole, TypeAccessToken, jwt.NewNumericDate(expiry))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, expiry, err
}

func GenerateRefreshToken(userID uint, businessId *uint, secret string, userRole usermodel.UserRole) (string, time.Time, error) {
	expiry := time.Now().Add(30 * 24 * time.Hour)
	claims := claimsGenerator(userID, businessId, secret, userRole, TypeRefreshToken, jwt.NewNumericDate(expiry))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, expiry, err
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

func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
