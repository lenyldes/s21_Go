package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolConfig() (*pgxpool.Config, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://postgres:postgres@localhost:5432/tictactoe?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 3

	return cfg, nil
}

func NewPostgresPool() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := NewPoolConfig()
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	queryCreateUserTable := "CREATE TABLE IF NOT EXISTS users (uuid TEXT PRIMARY KEY, login TEXT UNIQUE NOT NULL, password TEXT NOT NULL);"
	_, err = pool.Exec(ctx, queryCreateUserTable)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось создать таблицу users: %w", err)
	}

	queryCreateGamesTable := "CREATE TABLE IF NOT EXISTS games (uuid TEXT PRIMARY KEY, board TEXT[] NOT NULL, phase TEXT NOT NULL, player_uuid TEXT NOT NULL, player_x TEXT NOT NULL, player_o TEXT NOT NULL, created_at TIMESTAMP NOT NULL)"
	_, err = pool.Exec(ctx, queryCreateGamesTable)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось создать таблицу games: %w", err)
	}

	queryCreateRefreshTokensTable := "CREATE TABLE IF NOT EXISTS refresh_tokens (token_hash TEXT PRIMARY KEY, user_uuid TEXT NOT NULL REFERENCES users(uuid));"
	_, err = pool.Exec(ctx, queryCreateRefreshTokensTable)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось создать таблицу refresh_tokens: %w", err)
	}

	return pool, nil
}
