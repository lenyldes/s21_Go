package internal

import (
	"tictactoe/domain"
)

func gameToEntity(game domain.Game) (result GameEntity) {
	result.UUID = game.UUID
	result.CreatedAt = game.CreatedAt

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.Board[i*3+j] = string(game.Board[i][j])
		}
	}

	result.PlayerUUID = game.State.PlayerUUID
	result.PlayerX = game.PlayerX
	result.PlayerO = game.PlayerO
	result.Phase = string(game.State.Phase)

	return result
}

func gameToDomain(game GameEntity) (result domain.Game) {
	result.UUID = game.UUID
	result.CreatedAt = game.CreatedAt

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.Board[i][j] = domain.GameSymbol(game.Board[i*3+j])
		}
	}

	result.State.PlayerUUID = game.PlayerUUID
	result.PlayerX = game.PlayerX
	result.PlayerO = game.PlayerO
	result.State.Phase = domain.GamePhase(game.Phase)

	return result
}

func userToEntity(user domain.User) (result UserEntity) {
	result.UUID = user.UUID
	result.Login = user.Login
	result.Password = user.Password
	return result
}

func userToDomain(user UserEntity) (result domain.User) {
	result.UUID = user.UUID
	result.Login = user.Login
	result.Password = user.Password
	return result
}
