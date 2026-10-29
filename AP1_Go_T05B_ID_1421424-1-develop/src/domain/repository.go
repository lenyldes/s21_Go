package domain

import (
	"context"
)

type GameRepository interface {
	Save(ctx context.Context, game Game) error
	Get(ctx context.Context, UUID string) (Game, error)
	GetAvailableGames(ctx context.Context) ([]Game, error)
	GetCompletedGames(ctx context.Context, userUUID string) ([]Game, error)
	GetTopPlayers(ctx context.Context, n int) ([]WinStatistics, error)
	Create(ctx context.Context, game Game) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user User) (string, error) // возвращает UUID
	GetByUUID(ctx context.Context, UUID string) (User, error)
	GetByLogin(ctx context.Context, login string) (User, error)
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, userUUID string, token string) error
	Exists(ctx context.Context, token string) (bool, error)
	Replace(ctx context.Context, oldToken string, newToken string) (replaced bool, err error)
}
