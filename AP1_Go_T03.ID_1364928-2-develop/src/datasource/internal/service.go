package internal

import (
	"context"
	"errors"
	"tictactoe/domain"
)

type GameService struct {
	repository domain.GameRepository
}

func NewGameService(repository domain.GameRepository) *GameService {
	return &GameService{repository: repository}
}

func (s *GameService) IsFinished(ctx context.Context, game domain.Game) bool {
	// игра закончена если:
	// - компьютер выиграл
	// - человек выиграл
	// - нет свободных клеток (ничья)

	compWin := checkWinner(game.Board, domain.Computer)
	humanWin := checkWinner(game.Board, domain.Human)
	emptyCell := hasFreeCells(game.Board)

	return compWin || humanWin || !emptyCell
}

func checkWinner(board domain.Board, player int) bool {
	// проверяем 3 строки
	for i := 0; i < 3; i++ {
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true
		}
	}
	// проверяем 3 столбца
	for i := 0; i < 3; i++ {
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true
		}
	}
	// проверяем 2 диагонали
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[2][0] == player && board[1][1] == player && board[0][2] == player {
		return true
	}

	return false
}

func hasFreeCells(board domain.Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				return true
			}
		}
	}

	return false
}

func (s *GameService) Validate(ctx context.Context, savedGame domain.Game, newGame domain.Game) error {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if savedGame.Board[i][j] != domain.Empty && savedGame.Board[i][j] != newGame.Board[i][j] {
				return errors.New("клетка была изменена")
			}
		}
	}

	return nil
}

func (s *GameService) NextMove(ctx context.Context, game domain.Game) (domain.Game, error) {
	row, col := bestMove(game.Board)
	game.Board[row][col] = domain.Computer

	return game, nil
}

func minimax(board domain.Board, isComputerTurn bool) int {
	if checkWinner(board, domain.Computer) {
		return +1
	}
	if checkWinner(board, domain.Human) {
		return -1
	}
	if !hasFreeCells(board) {
		return 0
	}

	if isComputerTurn {
		bestScore := -10000

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domain.Empty {
					board[i][j] = domain.Computer
					score := minimax(board, !isComputerTurn)
					board[i][j] = domain.Empty
					if score > bestScore {
						bestScore = score
					}
				}

			}
		}
		return bestScore
	}

	bestScore := 10000
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				board[i][j] = domain.Human
				score := minimax(board, !isComputerTurn)
				board[i][j] = domain.Empty
				if score < bestScore {
					bestScore = score
				}
			}

		}
	}
	return bestScore
}

func bestMove(board domain.Board) (row, col int) {
	bestScore := -10000
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				board[i][j] = domain.Computer
				score := minimax(board, false)
				board[i][j] = domain.Empty
				if score > bestScore {
					bestScore = score
					row, col = i, j
				}
			}

		}
	}

	return row, col
}
