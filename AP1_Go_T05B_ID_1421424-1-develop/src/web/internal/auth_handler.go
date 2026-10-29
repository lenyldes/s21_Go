package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"tictactoe/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(s domain.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req domain.SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	userUUID, err := h.service.SignUp(r.Context(), req)

	if err != nil {
		if errors.Is(err, domain.ErrIncorrectRequest) {
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := UserResponse{UUID: userUUID}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при регистрации: %v", err)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req domain.JwtRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "некорректный запрос", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrIncorrectRequest) {
			http.Error(w, "некорректный запрос", http.StatusBadRequest)
			return
		}
		if errors.Is(err, domain.ErrIncorrectLoginOrPassword) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("ошибка записи ответа при входе: %v", err)
	}
}

func (h *AuthHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	h.refreshToken(w, r, h.service.RefreshAccessToken)
}

func (h *AuthHandler) RefreshRefreshToken(w http.ResponseWriter, r *http.Request) {
	h.refreshToken(w, r, h.service.RefreshRefreshToken)
}

func (h *AuthHandler) refreshToken(w http.ResponseWriter, r *http.Request, refresh func(context.Context, string) (domain.JwtResponse, error)) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req domain.RefreshJwtRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}
	if req.RefreshToken == "" {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	resp, err := refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenInvalidClaims) ||
			errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "Недействительный refreshToken", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Ошибка обновления токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("ошибка записи ответа при обновлении токена: %v", err)
	}
}
