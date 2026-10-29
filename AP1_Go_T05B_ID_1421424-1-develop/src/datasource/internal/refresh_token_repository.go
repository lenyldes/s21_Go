package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *RefreshTokenRepository) Save(ctx context.Context, userUUID string, token string) error {
	tokenHash := hashRefreshToken(token)
	query := "INSERT INTO refresh_tokens (token_hash, user_uuid) VALUES ($1, $2)"
	_, err := r.db.Exec(ctx, query, tokenHash, userUUID)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении refreshToken: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) Exists(ctx context.Context, token string) (bool, error) {
	tokenHash := hashRefreshToken(token)
	query := "SELECT EXISTS (SELECT 1 FROM refresh_tokens WHERE token_hash = $1)"
	var exists bool
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке наличия refreshToken: %w", err)
	}

	return exists, nil
}

func (r *RefreshTokenRepository) Replace(ctx context.Context, oldToken string, newToken string) (replaced bool, err error) {
	oldTokenHash := hashRefreshToken(oldToken)
	newTokenHash := hashRefreshToken(newToken)
	query := "UPDATE refresh_tokens SET token_hash = $1 WHERE token_hash = $2"
	result, err := r.db.Exec(ctx, query, newTokenHash, oldTokenHash)
	if err != nil {
		return false, fmt.Errorf("ошибка при замене refreshToken: %w", err)
	}

	return result.RowsAffected() == 1, nil
}
