package middleware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop/api"
	"shop/internal/web/dto"
)

type UserAuthenticator struct {
	client api.AuthServiceClient
}

func NewUserAuthenticator(client api.AuthServiceClient) *UserAuthenticator {
	return &UserAuthenticator{client: client}
}

func (a *UserAuthenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeAuthError(w, errors.New("некорректный заголовок авторизации"), http.StatusUnauthorized)
			return
		}

		_, err := a.client.ValidateToken(r.Context(), &api.AccessTokenRequest{
			AccessToken: parts[1], // В заголовке «Bearer <токен>» второй элемент — сам токен.
		})
		if err != nil {
			if status.Code(err) == codes.Unauthenticated || status.Code(err) == codes.InvalidArgument || status.Code(err) == codes.NotFound {
				writeAuthError(w, err, http.StatusUnauthorized)
				return
			}
			writeAuthError(w, err, http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeAuthError(w http.ResponseWriter, err error, status int) {
	response := dto.ErrorResponse{Errors: err.Error()}
	if status == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Bearer")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа с ошибкой аутентификации: %v", err)
	}
}
