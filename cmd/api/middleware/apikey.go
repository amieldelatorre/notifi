package middleware // import "github.com/amieldelatorre/notifi/cmd/api/middleware"

import "net/http"

func ApiKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do logic to check if Api Key is valid
		apiKey := r.Header.Get("ApiKey")
		if apiKey == "" || !isValidApiKey(apiKey) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// If valid set header for user-id
		r.Header.Set("x-user-id", "1234")

		next.ServeHTTP(w, r)
	})
}

func isValidApiKey(apiKey string) bool {
	return apiKey == "letmein"
}
