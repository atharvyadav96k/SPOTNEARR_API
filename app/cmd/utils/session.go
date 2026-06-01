package utils

import (
	"fmt"
	"time"

	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateSession(email string, expireTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(expireTime).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.C.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IsValidSession(session string) bool {
	token, err := jwt.Parse(session, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.C.JWTSecret), nil
	})
	return err == nil && token.Valid
}
