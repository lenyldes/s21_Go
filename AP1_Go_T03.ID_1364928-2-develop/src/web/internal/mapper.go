package internal

import (
	"tictactoe/domain"
)

func toDomain(req GameRequest, UUID string) (result domain.Game) {
	result.UUID = UUID
	result.Board = domain.Board(req.Board)
	return result
}

func toWeb(game domain.Game, gameOver bool) (result GameResponse) {
	result.Board = BoardWeb(game.Board)
	result.UUID = game.UUID
	result.GameOver = gameOver
	return result
}
