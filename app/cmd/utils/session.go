package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateSession(email string, expireTime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(expireTime).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
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
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return err == nil && token.Valid
}
