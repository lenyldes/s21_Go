package internal

import (
	"context"
	"tictactoe/domain"
)

type GameRepository struct {
	storage Storage
}

func NewGameRepository(storage Storage) *GameRepository {
	return &GameRepository{storage: storage}
}

func (r *GameRepository) Save(ctx context.Context, game domain.Game) error {
	return r.storage.Save(ctx, game)
}

func (r *GameRepository) Get(ctx context.Context, UUID string) (result domain.Game, err error) {
	return r.storage.Get(ctx, UUID)
}

func (r *GameRepository) Create(ctx context.Context, game domain.Game) (domain.Game, error) {
	return r.storage.Create(ctx, game)
}
