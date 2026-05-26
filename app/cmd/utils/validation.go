package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func ValidateEmail(email string) bool {
	cleanEmail := strings.TrimSpace(email)
	if cleanEmail == "" {
		return false
	}
	return emailRegex.MatchString(cleanEmail)
}

func ValidatePhone(phone string) bool {
	cleanPhone := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(phone)
	if cleanPhone == "" {
		return false
	}
	return phoneRegex.MatchString(cleanPhone)
}
func ValidatePassword(password string) (bool, string) {
	var (
		hasMinLen = len(password) >= 8
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasMinLen {
		return false, "Password must be at least 8 characters long"
	}
	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one numeric digit"
	}

	return true, ""
}
