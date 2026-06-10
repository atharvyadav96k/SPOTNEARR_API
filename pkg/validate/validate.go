package validate

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func Email(email string) bool {
	clean := strings.TrimSpace(email)
	if clean == "" {
		return false
	}
	return emailRegex.MatchString(clean)
}

func Phone(phone string) bool {
	clean := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(phone)
	if clean == "" {
		return false
	}
	return phoneRegex.MatchString(clean)
}

func Password(password string) (bool, string) {
	var hasUpper, hasLower, hasNumber bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}
	if len(password) < 8 {
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
