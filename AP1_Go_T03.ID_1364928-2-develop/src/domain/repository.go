package domain

import "context"

type GameRepository interface {
	Save(ctx context.Context, game Game) error
	Get(ctx context.Context, UUID string) (Game, error)
	Create(ctx context.Context, game Game) (Game, error)
}
