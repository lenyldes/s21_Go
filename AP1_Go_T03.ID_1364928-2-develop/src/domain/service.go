package domain

import "context"

type GameService interface {
	NextMove(ctx context.Context, game Game) (Game, error)
	Validate(ctx context.Context, savedGame Game, newGame Game) error
	IsFinished(ctx context.Context, game Game) bool
}
