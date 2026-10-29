package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop/api"
	"shop/internal/web/dto"
)

type AuthHandler struct {
	client api.AuthServiceClient
}

func NewAuthHandler(client api.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client: client}
}

// @Summary Register user
// @Description Доступен без токена. Создаёт пользователя через сервис авторизации и возвращает access_token и refresh_token. Поля login, email и role не могут быть пустыми или состоять только из пробелов; password должен быть непустым. Поля name, surname и phonenumber необязательны и не валидируются магазином. Значения передаются без обрезки пробелов. Проверку требований к паролю и хэширование выполняет сервис авторизации.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "Registration data"
// @Success 201 {object} dto.SignupResponse "User created"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON, required fields or registration data rejected by the auth service"
// @Failure 409 {object} dto.ErrorResponse "Login already exists"
// @Failure 500 {object} dto.ErrorResponse "Failed to register user, including auth service unavailability"
// @Router /register [post]
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var request dto.SignupRequest
	err := decodeJSONBody(r.Body, &request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	tokens, err := h.client.Signup(r.Context(), &api.SignupRequest{
		Login:       request.Login,
		Email:       request.Email,
		Role:        request.Role,
		Password:    request.Password,
		Name:        request.Name,
		Surname:     request.Surname,
		Phonenumber: request.Phonenumber,
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			writeError(w, err, http.StatusConflict)
			return
		}
		if status.Code(err) == codes.InvalidArgument {
			writeError(w, err, http.StatusBadRequest)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := dto.SignupResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при регистрации: %v", err)
	}
}

// @Summary Sign in
// @Description Доступен без токена. Проверяет логин и пароль через сервис авторизации и возвращает access_token и refresh_token. Поле login не может быть пустым или состоять только из пробелов; password должен быть непустым. Значения передаются без изменений. Для доступа к защищённым методам магазина нужно передавать access_token в заголовке Authorization: Bearer <access_token>.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SigninRequest true "Sign-in data"
// @Success 200 {object} dto.SigninResponse "Access and refresh tokens"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON or required fields"
// @Failure 401 {object} dto.ErrorResponse "User not found or invalid password"
// @Failure 500 {object} dto.ErrorResponse "Failed to sign in, including auth service unavailability"
// @Router /auth [post]
func (h *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var request dto.SigninRequest
	err := decodeJSONBody(r.Body, &request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	tokens, err := h.client.Signin(r.Context(), &api.SigninRequest{
		Login:    request.Login,
		Password: request.Password,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound || status.Code(err) == codes.Unauthenticated {
			writeError(w, err, http.StatusUnauthorized)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := dto.SigninResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при входе: %v", err)
	}
}

// @Summary Reset password
// @Description Доступен без токена. Передаёт логин в сервис авторизации, который заменяет пароль временным и выводит его в консоль вместо отправки письма. Поле login не может быть пустым или состоять только из пробелов; значение передаётся без обрезки пробелов. У пользователя должен быть указан email. Магазин не получает и не возвращает временный пароль. Успешный ответ не содержит тела.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Password reset data"
// @Success 204 "Password reset"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON, login or missing user email"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to reset password, including auth service unavailability"
// @Router /reset [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var request dto.ResetPasswordRequest
	err := decodeJSONBody(r.Body, &request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	_, err = h.client.ResetPassword(r.Context(), &api.ResetPasswordRequest{
		Login: request.Login,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			writeError(w, err, http.StatusNotFound)
			return
		}
		if status.Code(err) == codes.InvalidArgument {
			writeError(w, err, http.StatusBadRequest)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Refresh
// @Description Обновляет refresh токен
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "request"
// @Success 200 {object} dto.RefreshTokensResponse "Access and refresh tokens"
// @Failure 400 {object} dto.ErrorResponse "Invalid refresh token"
// @Failure 401 {object} dto.ErrorResponse "User not found or invalid refresh token"
// @Failure 500 {object} dto.ErrorResponse "Failed to refresh, including auth service unavailability"
// @Router /refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request dto.RefreshTokenRequest
	err := decodeJSONBody(r.Body, &request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	tokens, err := h.client.Refresh(r.Context(), &api.RefreshTokensRequest{
		RefreshToken: *request.RefreshToken,
	})
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			writeError(w, err, http.StatusUnauthorized)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	response := dto.RefreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при обновении токена: %v", err)
	}
}

// @Summary ChangePassword
// @Description Изменяет пароль пользователя
// @Tags Auth
// @Accept json
// @Param request body dto.ChangePasswordRequest true "request"
// @Success 204 "Password reset"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON, login or missing user email"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to reset password, including auth service unavailability"
// @Router /change_password [post]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var request dto.ChangePasswordRequest
	err := decodeJSONBody(r.Body, &request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	_, err = h.client.ChangePassword(r.Context(), &api.ChangePasswordRequest{
		Login:           *request.Login,
		CurrentPassword: *request.CurrentPassword,
		NewPassword:     *request.NewPassword,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			writeError(w, err, http.StatusNotFound)
			return
		}
		if status.Code(err) == codes.Unauthenticated {
			writeError(w, err, http.StatusUnauthorized)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
