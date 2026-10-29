package auth

import (
	"auth/internal/shared"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken string
type RefreshToken string

const (
	REFRESH_TOKEN_TYPE = "refresh"
	ACCESS_TOKEN_TYPE  = "access"
)

type JwtProvider struct {
	Options         []jwt.ParserOption
	keys            keys
	JwtCache        *JwtCache
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

type keys struct {
	accessKey  []byte
	refreshKey []byte
}

func NewKeys(config *shared.Config) keys {
	return keys{
		accessKey:  []byte(config.AccessTokenKey),
		refreshKey: []byte(config.RefreshTokenKey),
	}
}

func NewJwtProvider(options []jwt.ParserOption, keys keys, jwtCache *JwtCache, config *shared.Config) *JwtProvider {
	return &JwtProvider{
		Options:         options,
		keys:            keys,
		JwtCache:        jwtCache,
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
		issuer:          config.Issuer,
	}
}

func (p *JwtProvider) IssueTokens(userId uuid.UUID, role string) (AccessToken, RefreshToken, error) {
	accessToken, err := p.issueAccessToken(userId, role)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := p.issueRefreshTokenForNewFamily(userId, role)
	if err != nil {
		return "", "", err
	}
	return accessToken, RefreshToken(refreshToken), nil
}

func (p *JwtProvider) issueAccessToken(userId uuid.UUID, role string) (AccessToken, error) {
	accessClaims := AccessClaims{
		UserId:    userId.String(),
		Role:      role,
		TokenType: ACCESS_TOKEN_TYPE,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.accessTokenTTL)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(p.keys.accessKey)
	if err != nil {
		return "", err
	}
	return AccessToken(accessTokenString), nil
}

func (p *JwtProvider) issueRefreshTokenForNewFamily(userId uuid.UUID, role string) (RefreshToken, error) {
	tokenId := p.JwtCache.CreateFamily(userId, role)
	refreshClaims := RefreshClaims{
		TokenType: REFRESH_TOKEN_TYPE,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Issuer:    p.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTokenTTL)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(p.keys.refreshKey)
	if err != nil {
		return "", err
	}
	return RefreshToken(refreshTokenString), nil
}

func (p *JwtProvider) issueRefreshTokenForExistsFamily(cur uuid.UUID) (RefreshToken, error) {
	tokenId, err := p.JwtCache.Rotate(cur)
	if err != nil {
		return "", err
	}
	refreshClaims := RefreshClaims{
		TokenType: REFRESH_TOKEN_TYPE,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Issuer:    p.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTokenTTL)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(p.keys.refreshKey)
	if err != nil {
		return "", err
	}
	return RefreshToken(refreshTokenString), nil
}

func (p *JwtProvider) Refresh(refreshToken string) (AccessToken, RefreshToken, string, error) {
	var refreshClaims RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, p.refreshTokenKeyFunc, p.Options...)
	if err != nil || !token.Valid {
		if err == nil {
			err = fmt.Errorf("token invalid")
		}
		return "", "", "", err
	}
	if refreshClaims.TokenType != REFRESH_TOKEN_TYPE || refreshClaims.ID == "" {
		return "", "", "", errors.New("not refresh token")
	}

	cur, err := uuid.Parse(refreshClaims.ID)
	newRefreshToken, err := p.issueRefreshTokenForExistsFamily(cur)
	if err != nil {
		return "", "", "", err
	}
	userId, role, err := p.JwtCache.LookUpUser(cur)
	if err != nil {
		return "", "", "", err
	}
	accessToken, err := p.issueAccessToken(userId, role)
	if err != nil {
		return "", "", "", err
	}
	return accessToken, newRefreshToken, userId.String(), nil
}

func (p *JwtProvider) AccessTokenKeyFunc(_ *jwt.Token) (any, error) {
	return p.keys.accessKey, nil
}

func (p *JwtProvider) refreshTokenKeyFunc(_ *jwt.Token) (any, error) {
	return p.keys.refreshKey, nil
}

func Options(config *shared.Config) []jwt.ParserOption {
	return []jwt.ParserOption{
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(config.Leeway),
		jwt.WithIssuer(config.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	}
}
