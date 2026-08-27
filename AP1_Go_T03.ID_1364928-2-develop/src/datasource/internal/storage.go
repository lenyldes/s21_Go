package internal

import (
	"context"
	"errors"
	"sync"
	"tictactoe/domain"
)

var ErrGameNotFound = errors.New("игра не найдена")

type GameStorage struct {
	storage sync.Map
}

func NewGameStorage() *GameStorage {
	return &GameStorage{}
}

func (storage *GameStorage) Save(ctx context.Context, game domain.Game) error {
	gameEntity := toEntity(game)
	storage.storage.Store(game.UUID, gameEntity)

	return nil
}

func (storage *GameStorage) Get(ctx context.Context, UUID string) (result domain.Game, err error) {
	value, ok := storage.storage.Load(UUID)
	if !ok {
		err = ErrGameNotFound
		return result, err
	}

	entity := value.(GameEntity)
	result = toDomain(entity)

	return result, err
}

func (storage *GameStorage) Create(ctx context.Context, game domain.Game) (domain.Game, error) {
	err := storage.Save(ctx, game)

	return game, err
}
