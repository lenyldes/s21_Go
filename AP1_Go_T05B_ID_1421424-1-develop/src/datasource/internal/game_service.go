package internal

import (
	"context"
	"fmt"
	"tictactoe/domain"
	"time"

	"github.com/google/uuid"
)

type GameService struct {
	repository domain.GameRepository
}

func NewGameService(repository domain.GameRepository) *GameService {
	return &GameService{repository: repository}
}

func (s *GameService) GetCompletedGames(ctx context.Context, userUUID string) ([]domain.Game, error) {
	return s.repository.GetCompletedGames(ctx, userUUID)
}

func (s *GameService) GetTopPlayers(ctx context.Context, n int) ([]domain.WinStatistics, error) {
	return s.repository.GetTopPlayers(ctx, n)
}

func (s *GameService) CreateGame(ctx context.Context, mode domain.GameMode, userUUID string) (domain.Game, error) {
	if mode != domain.GameModeVSComputer && mode != domain.GameModeVSHuman {
		return domain.Game{}, domain.ErrIncorrectGameMode
	}

	var newGame domain.Game

	newGame.UUID = uuid.New().String()
	newGame.CreatedAt = time.Now().UTC()
	newGame.PlayerX = userUUID

	if mode == domain.GameModeVSComputer {
		newGame.State.PlayerUUID = userUUID
		newGame.PlayerO = domain.ComputerUUID
		newGame.State.Phase = domain.GamePhaseMove
	}
	if mode == domain.GameModeVSHuman {
		newGame.State.Phase = domain.GamePhaseReady
	}

	_, err := s.repository.Create(ctx, newGame)
	if err != nil {
		return domain.Game{}, fmt.Errorf("ошибка при создании игры: %w", err)
	}

	return newGame, nil
}

func (s *GameService) JoinGame(ctx context.Context, gameUUID string, userUUID string) (domain.Game, error) {
	game, err := s.repository.Get(ctx, gameUUID)
	if err != nil {
		return domain.Game{}, fmt.Errorf("ошибка при присоединении к игре: %w", err)
	}
	if game.PlayerX == userUUID || game.PlayerO == userUUID {
		return domain.Game{}, domain.ErrAlreadyInGame
	}
	if game.State.Phase != domain.GamePhaseReady || game.PlayerO != "" {
		return domain.Game{}, domain.ErrGameUnavailable
	}

	game.PlayerO = userUUID
	game.State.Phase = domain.GamePhaseMove
	game.State.PlayerUUID = game.PlayerX

	err = s.repository.Save(ctx, game)
	if err != nil {
		return domain.Game{}, fmt.Errorf("ошибка при сохранении игры: %w", err)
	}

	return game, nil
}

func (s *GameService) IsFinished(ctx context.Context, game domain.Game) bool {
	if game.State.Phase == domain.GamePhaseWin || game.State.Phase == domain.GamePhaseDraw {
		return true
	}

	return false
}

func updateGameState(game domain.Game) domain.Game {
	playerXWin := checkWinner(game.Board, domain.GameSymbolX)
	if playerXWin == true {
		game.State.PlayerUUID = game.PlayerX
		game.State.Phase = domain.GamePhaseWin
		return game
	}

	playerOWin := checkWinner(game.Board, domain.GameSymbolO)
	if playerOWin == true {
		game.State.PlayerUUID = game.PlayerO
		game.State.Phase = domain.GamePhaseWin
		return game
	}

	emptyCell := hasFreeCells(game.Board)
	if !emptyCell {
		game.State.PlayerUUID = ""
		game.State.Phase = domain.GamePhaseDraw
		return game
	}

	if game.State.PlayerUUID == game.PlayerX {
		game.State.PlayerUUID = game.PlayerO
	} else {
		game.State.PlayerUUID = game.PlayerX
	}

	game.State.Phase = domain.GamePhaseMove

	return game
}

func checkWinner(board domain.Board, symbol domain.GameSymbol) bool {
	// проверяем 3 строки
	for i := 0; i < 3; i++ {
		if board[i][0] == symbol && board[i][1] == symbol && board[i][2] == symbol {
			return true
		}
	}
	// проверяем 3 столбца
	for i := 0; i < 3; i++ {
		if board[0][i] == symbol && board[1][i] == symbol && board[2][i] == symbol {
			return true
		}
	}
	// проверяем 2 диагонали
	if board[0][0] == symbol && board[1][1] == symbol && board[2][2] == symbol {
		return true
	}
	if board[2][0] == symbol && board[1][1] == symbol && board[0][2] == symbol {
		return true
	}

	return false
}

func hasFreeCells(board domain.Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.GameSymbolEmpty {
				return true
			}
		}
	}

	return false
}

func (s *GameService) Validate(ctx context.Context, savedGame domain.Game, newGame domain.Game, userUUID string) error {
	if savedGame.State.Phase != domain.GamePhaseMove || savedGame.State.PlayerUUID != userUUID {
		return domain.ErrMoveNotAllowed
	}

	if savedGame.PlayerX != userUUID && savedGame.PlayerO != userUUID {
		return domain.ErrMoveNotAllowed
	}

	var allowedSymbol domain.GameSymbol
	if userUUID == savedGame.PlayerX {
		allowedSymbol = domain.GameSymbolX
	} else {
		allowedSymbol = domain.GameSymbolO
	}

	var countChangeCell int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if savedGame.Board[i][j] != newGame.Board[i][j] {
				countChangeCell++
				// можем изменить только одну клетку
				if countChangeCell > 1 {
					return domain.ErrInvalidMove
				}

				// можем изменять только пустую клетку
				if savedGame.Board[i][j] != domain.GameSymbolEmpty {
					return domain.ErrInvalidMove
				}

				if newGame.Board[i][j] != allowedSymbol {
					return domain.ErrInvalidMove
				}
			}
		}
	}
	if countChangeCell != 1 {
		return domain.ErrInvalidMove
	}

	return nil
}

func (s *GameService) ApplyHumanMove(ctx context.Context, savedGame domain.Game, newGame domain.Game, userUUID string) (domain.Game, error) {
	err := s.Validate(ctx, savedGame, newGame, userUUID)

	if err != nil {
		return domain.Game{}, err
	}

	game := savedGame
	game.Board = newGame.Board
	game = updateGameState(game)

	return game, nil
}

func (s *GameService) NextMove(ctx context.Context, game domain.Game) (domain.Game, error) {
	// за X играет человек
	if game.PlayerX == "" || game.PlayerX == domain.ComputerUUID {
		return domain.Game{}, domain.ErrMoveNotAllowed
	}

	// за O играет компьютер
	if game.PlayerO != domain.ComputerUUID {
		return domain.Game{}, domain.ErrMoveNotAllowed
	}

	// игра в фазе хода
	if game.State.Phase != domain.GamePhaseMove {
		return domain.Game{}, domain.ErrMoveNotAllowed
	}

	// сейчас очередь компьютера
	if game.State.PlayerUUID != domain.ComputerUUID {
		return domain.Game{}, domain.ErrMoveNotAllowed
	}

	row, col := bestMove(game.Board)
	game.Board[row][col] = domain.GameSymbolO
	game = updateGameState(game)

	return game, nil
}

func minimax(board domain.Board, isComputerTurn bool) int {
	// человек - X: domain.GameSymbolX
	// компьютер - O: domain.GameSymbolO
	if checkWinner(board, domain.GameSymbolO) {
		return +1
	}
	if checkWinner(board, domain.GameSymbolX) {
		return -1
	}
	if !hasFreeCells(board) {
		return 0
	}

	if isComputerTurn {
		bestScore := -10000

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domain.GameSymbolEmpty {
					board[i][j] = domain.GameSymbolO
					score := minimax(board, !isComputerTurn)
					board[i][j] = domain.GameSymbolEmpty
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
			if board[i][j] == domain.GameSymbolEmpty {
				board[i][j] = domain.GameSymbolX
				score := minimax(board, !isComputerTurn)
				board[i][j] = domain.GameSymbolEmpty
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
			if board[i][j] == domain.GameSymbolEmpty {
				board[i][j] = domain.GameSymbolO
				score := minimax(board, false)
				board[i][j] = domain.GameSymbolEmpty
				if score > bestScore {
					bestScore = score
					row, col = i, j
				}
			}

		}
	}

	return row, col
}
