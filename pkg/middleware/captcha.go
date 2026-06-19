package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type turnstileResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

// CaptchaValidation verifies a Cloudflare Turnstile token from X-Captcha-Token.
func CaptchaValidation(captchaURL, captchaSecret string) func(http.Handler) http.Handler {
	client := &http.Client{Timeout: 3 * time.Second}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captchaToken := r.Header.Get("X-Captcha-Token")
			if captchaToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message": "Unauthorized access: X-Captcha-Token header is missing"}`))
				return
			}
			formData := url.Values{}
			formData.Set("secret", captchaSecret)
			formData.Set("response", captchaToken)

			resp, err := client.PostForm(captchaURL, formData)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal verification system error"}`))
				return
			}
			defer resp.Body.Close()

			var cfResp turnstileResponse
			if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Error reading security verification response"}`))
				return
			}
			if !cfResp.Success {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"message": "Security validation failed. Automated traffic blocked.",
					"errors":  cfResp.ErrorCodes,
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
