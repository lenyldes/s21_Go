package internal

import (
	"context"
	"tictactoe/domain"
)

type BoardEntity []int

type GameEntity struct {
	UUID  string
	Board BoardEntity
}

type Storage interface {
	Save(ctx context.Context, game domain.Game) error
	Get(ctx context.Context, UUID string) (domain.Game, error)
	Create(ctx context.Context, game domain.Game) (domain.Game, error)
}
