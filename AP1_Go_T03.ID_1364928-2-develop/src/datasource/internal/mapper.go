package internal

import (
	"tictactoe/domain"
)

func toEntity(game domain.Game) (result GameEntity) {
	result.Board = make(BoardEntity, 9)

	result.UUID = game.UUID

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.Board[i*3+j] = game.Board[i][j]
		}
	}

	return result
}

func toDomain(game GameEntity) (result domain.Game) {
	result.UUID = game.UUID

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.Board[i][j] = game.Board[i*3+j]
		}
	}

	return result
}
