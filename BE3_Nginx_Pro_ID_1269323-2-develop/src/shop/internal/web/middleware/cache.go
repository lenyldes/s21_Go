package middleware

import "net/http"

// NoStore помечает ответы приложения как не предназначенные для хранения в HTTP-кешах.
func NoStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
