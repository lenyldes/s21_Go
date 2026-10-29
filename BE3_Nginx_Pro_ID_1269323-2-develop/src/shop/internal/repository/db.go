package repository

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
		url = "postgres://postgres:postgres@localhost:5432/shop?sslmode=disable"
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

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	return pool, nil
}
