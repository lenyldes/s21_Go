package internal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"tictactoe/domain"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userUUID := r.PathValue("id")
	user, err := h.service.GetByUUID(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, domain.ErrUserNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UserInfoResponse{
		UUID:  user.UUID,
		Login: user.Login,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи информации о пользователе: %v", err)
	}
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	value := r.Context().Value(userUUIDKey)
	userUUID, ok := value.(string)
	if userUUID == "" || !ok {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetByUUID(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, domain.ErrUserNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UserInfoResponse{
		UUID:  user.UUID,
		Login: user.Login,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи информации о текущем пользователе: %v", err)
	}
}
