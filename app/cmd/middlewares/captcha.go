package middleware

import (
	"net/http"
)

func CaptchaValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// development bypass
		next.ServeHTTP(w, r)
		return
		captchaHeader := r.Header.Get("X-Captcha-Token")
		if captchaHeader == "" {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		// captcha validation

		next.ServeHTTP(w, r)

	})
}
