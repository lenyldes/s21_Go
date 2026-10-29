package internal

import (
	"context"
	"net/http"
	"strings"
	"tictactoe/domain"
)

type UserAuthenticator struct {
	provider domain.JwtProvider
}

func NewUserAuthenticator(provider domain.JwtProvider) *UserAuthenticator {
	return &UserAuthenticator{
		provider: provider,
	}
}

func (a *UserAuthenticator) Authenticate(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "Некорректный заголовок авторизации", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		err := a.provider.ValidateAccessToken(token)
		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		userUUID, err := a.provider.GetUUIDFromToken(token)
		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userUUIDKey, userUUID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(handler)
}
