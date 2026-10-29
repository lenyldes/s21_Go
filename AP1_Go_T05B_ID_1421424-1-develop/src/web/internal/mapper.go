package internal

import (
	"tictactoe/domain"
)

func toDomain(req GameRequest, UUID string) (result domain.Game) {
	result.UUID = UUID
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.Board[i][j] = domain.GameSymbol(req.Board[i][j])
		}
	}
	return result
}

func toWeb(game domain.Game, gameOver bool) (result GameResponse) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.Board[i][j] = string(game.Board[i][j])
		}
	}

	result.UUID = game.UUID
	result.CreatedAt = game.CreatedAt
	result.GameOver = gameOver

	result.Phase = string(game.State.Phase)
	result.PlayerUUID = game.State.PlayerUUID
	result.PlayerX = game.PlayerX
	result.PlayerO = game.PlayerO
	return result
}
