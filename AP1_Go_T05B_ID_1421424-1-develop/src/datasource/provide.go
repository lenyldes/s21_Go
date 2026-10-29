package datasource

import (
	"tictactoe/datasource/internal"
	"tictactoe/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewGameRepository(db *pgxpool.Pool) domain.GameRepository {
	return internal.NewGameRepository(db)
}

func NewGameService(repository domain.GameRepository) domain.GameService {
	return internal.NewGameService(repository)
}

func NewPostgresPool() (*pgxpool.Pool, error) {
	return internal.NewPostgresPool()
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return internal.NewUserRepository(db)
}

func NewRefreshTokenRepository(db *pgxpool.Pool) domain.RefreshTokenRepository {
	return internal.NewRefreshTokenRepository(db)
}

func NewUserService(repository domain.UserRepository) domain.UserService {
	return internal.NewUserService(repository)
}

func NewJwtProvider() domain.JwtProvider {
	return internal.NewJwtProvider()
}

func NewAuthService(service domain.UserService, provider domain.JwtProvider, repository domain.RefreshTokenRepository) domain.AuthService {
	return internal.NewAuthService(service, provider, repository)
}
