package internal

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	service    domain.UserService
	provider   domain.JwtProvider
	repository domain.RefreshTokenRepository
}

func NewAuthService(s domain.UserService, provider domain.JwtProvider, repository domain.RefreshTokenRepository) *AuthService {
	return &AuthService{
		service:    s,
		provider:   provider,
		repository: repository,
	}
}

func (a *AuthService) SignUp(ctx context.Context, req domain.SignUpRequest) (string, error) {
	if req.Login == "" || req.Password == "" {
		return "", domain.ErrIncorrectRequest
	}

	newUUID := uuid.New().String()

	user := domain.User{
		UUID:     newUUID,
		Login:    req.Login,
		Password: req.Password,
	}

	userUUID, err := a.service.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	return userUUID, nil
}

func (a *AuthService) Login(ctx context.Context, req domain.JwtRequest) (domain.JwtResponse, error) { // возвращает токены авторизации
	if req.Login == "" || req.Password == "" {
		return domain.JwtResponse{}, fmt.Errorf("логин и пароль обязательны: %w", domain.ErrIncorrectRequest)
	}

	user, err := a.service.GetByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.JwtResponse{}, fmt.Errorf("ошибка при входе: %w", domain.ErrIncorrectLoginOrPassword)
		}
		return domain.JwtResponse{}, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	if user.Password != req.Password {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при входе: %w", domain.ErrIncorrectLoginOrPassword)
	}

	accessToken, err := a.provider.GenerateAccessToken(user)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при создании accessToken: %w", err)
	}

	refreshToken, err := a.provider.GenerateRefreshToken(user)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при создании refreshToken: %w", err)
	}

	err = a.repository.Save(ctx, user.UUID, refreshToken)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при сохранении refreshToken пользователя: %w", err)
	}

	return domain.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) RefreshRefreshToken(ctx context.Context, refreshToken string) (domain.JwtResponse, error) {
	err := a.provider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при проверке refreshToken: %w: %w", jwt.ErrTokenInvalidClaims, err)
	}

	userUUID, err := a.provider.GetUUIDFromToken(refreshToken)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при получении UUID из refreshToken: %w", err)
	}

	user, err := a.service.GetByUUID(ctx, userUUID)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	accessToken, err := a.provider.GenerateAccessToken(user)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при создании accessToken: %w", err)
	}

	newRefreshToken, err := a.provider.GenerateRefreshToken(user)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при создании refreshToken: %w", err)
	}

	replaced, err := a.repository.Replace(ctx, refreshToken, newRefreshToken)
	if err != nil {
		return domain.JwtResponse{}, fmt.Errorf("ошибка при замене refreshToken пользователя: %w", err)
	}
	if !replaced {
		return domain.JwtResponse{}, fmt.Errorf("refreshToken отсутствует или уже заменён: %w", jwt.ErrTokenInvalidClaims)
	}

	return domain.JwtResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (domain.JwtResponse, error) {
	return a.RefreshRefreshToken(ctx, refreshToken)
}
