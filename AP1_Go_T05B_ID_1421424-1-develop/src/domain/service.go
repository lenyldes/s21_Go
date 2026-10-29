package domain

import "context"

type GameService interface {
	ApplyHumanMove(ctx context.Context, savedGame Game, newGame Game, userUUID string) (Game, error)
	NextMove(ctx context.Context, game Game) (Game, error)
	Validate(ctx context.Context, savedGame Game, newGame Game, userUUID string) error
	IsFinished(ctx context.Context, game Game) bool
	CreateGame(ctx context.Context, mode GameMode, userUUID string) (Game, error)
	JoinGame(ctx context.Context, gameUUID string, userUUID string) (Game, error)
	GetCompletedGames(ctx context.Context, userUUID string) ([]Game, error)
	GetTopPlayers(ctx context.Context, n int) ([]WinStatistics, error)
}

type UserService interface {
	Create(ctx context.Context, user User) (string, error) // возвращает UUID
	GetByUUID(ctx context.Context, UUID string) (User, error)
	GetByLogin(ctx context.Context, login string) (User, error)
}

type JwtProvider interface {
	GenerateAccessToken(user User) (string, error)
	GenerateRefreshToken(user User) (string, error)
	ValidateAccessToken(token string) error
	ValidateRefreshToken(token string) error
	GetUUIDFromToken(token string) (string, error)
}

type AuthService interface {
	SignUp(ctx context.Context, req SignUpRequest) (string, error)  // возвращает UUID созданного пользователя
	Login(ctx context.Context, req JwtRequest) (JwtResponse, error) // возвращает токены авторизации
	RefreshAccessToken(ctx context.Context, refreshToken string) (JwtResponse, error)
	RefreshRefreshToken(ctx context.Context, refreshToken string) (JwtResponse, error)
}
