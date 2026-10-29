package domain

import (
	"errors"
	"time"
)

var ErrGameNotFound = errors.New("игра не найдена")
var ErrUserNotFound = errors.New("пользователь не найден")
var ErrIncorrectLoginOrPassword = errors.New("некорректный логин или пароль")
var ErrIncorrectRequest = errors.New("некорректный формат заголовка/токена")
var ErrIncorrectGameMode = errors.New("недопустимый режим игры")
var ErrGameUnavailable = errors.New("игра недоступна")
var ErrAlreadyInGame = errors.New("пользователь уже участвует в игре")
var ErrMoveNotAllowed = errors.New("ход запрещен")
var ErrInvalidMove = errors.New("поле изменено неверно")

const ComputerUUID string = "00000000-0000-0000-0000-000000000000"

type GameSymbol string

const (
	GameSymbolX     GameSymbol = "X"
	GameSymbolO     GameSymbol = "O"
	GameSymbolEmpty GameSymbol = ""
)

type Board [3][3]GameSymbol

type Game struct {
	UUID      string // UUID игры
	Board     Board
	State     GameState
	PlayerX   string // UUID участника с X: ComputerUUID — компьютер, "" — место свободно
	PlayerO   string // UUID участника с O: ComputerUUID — компьютер, "" — место свободно
	CreatedAt time.Time
}

type User struct {
	UUID     string
	Login    string
	Password string
}

type WinStatistics struct {
	UserUUID string
	WinRatio float64
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type GamePhase string

const (
	GamePhaseReady GamePhase = "ready" // ожидание игроков
	GamePhaseMove  GamePhase = "move"  // игра в процессе, ход игрока
	GamePhaseWin   GamePhase = "win"   // окончание игры, победа одного из игроков
	GamePhaseDraw  GamePhase = "draw"  // окончение игры, ничья
)

type GameState struct {
	PlayerUUID string // UUID текущего игрока в фазе move или победителя в фазе win; "" в ready и draw
	Phase      GamePhase
}

type GameMode string

const (
	GameModeVSHuman    GameMode = "vs_human"
	GameModeVSComputer GameMode = "vs_computer"
)
