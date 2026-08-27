package datasource

import (
	"tictactoe/datasource/internal"
	"tictactoe/domain"
)

func NewGameStorage() *internal.GameStorage {
	return internal.NewGameStorage()
}

func NewGameRepository(storage *internal.GameStorage) domain.GameRepository {
	return internal.NewGameRepository(storage)
}

func NewGameService(repository domain.GameRepository) domain.GameService {
	return internal.NewGameService(repository)
}
