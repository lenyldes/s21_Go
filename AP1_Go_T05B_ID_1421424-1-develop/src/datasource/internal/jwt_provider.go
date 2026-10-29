package internal

import (
	"fmt"
	"os"
	"tictactoe/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	accessTokenLifetime  = 15 * time.Minute
	refreshTokenLifetime = 7 * 24 * time.Hour
)

type JwtProvider struct {
	secret []byte
}

func NewJwtProvider() *JwtProvider {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // secret - учебная заглушка, в проде конечно же так делать не будем)
	}

	return &JwtProvider{secret: []byte(secret)}
}

func (p *JwtProvider) GenerateAccessToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":        user.UUID,
		"exp":        time.Now().Add(accessTokenLifetime).Unix(),
		"token_type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *JwtProvider) GenerateRefreshToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":        user.UUID,
		"exp":        time.Now().Add(refreshTokenLifetime).Unix(),
		"jti":        uuid.New().String(),
		"token_type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *JwtProvider) ValidateAccessToken(token string) error {
	return p.validateToken(token, "access")
}

func (p *JwtProvider) ValidateRefreshToken(token string) error {
	return p.validateToken(token, "refresh")
}

func (p *JwtProvider) GetUUIDFromToken(token string) (string, error) {
	_, userUUID, err := p.parseToken(token)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении UUID из токена: %w", err)
	}

	return userUUID, nil
}

func (p *JwtProvider) validateToken(token string, tokenType string) error {
	claims, _, err := p.parseToken(token)
	if err != nil {
		return fmt.Errorf("ошибка при валидации %s-токена: %w", tokenType, err)
	}

	if claims["token_type"] != tokenType {
		return fmt.Errorf("неверный тип токена: %w", jwt.ErrTokenInvalidClaims)
	}

	return nil
}

func (p *JwtProvider) parseToken(token string) (claims jwt.MapClaims, userUUID string, err error) {
	keyFunc := func(_ *jwt.Token) (any, error) {
		return p.secret, nil
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	claims = jwt.MapClaims{}

	_, err = parser.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при проверке токена: %w", err)
	}

	userUUID, err = claims.GetSubject()
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при чтении UUID из токена: %w", err)
	}
	_, err = uuid.Parse(userUUID)
	if err != nil {
		return nil, "", fmt.Errorf("некорректный UUID пользователя в токене: %w", err)
	}

	return claims, userUUID, nil
}
