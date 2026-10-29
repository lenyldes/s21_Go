package domain

import "context"

// В методе List int64 соответствует диапазону LIMIT/OFFSET PostgreSQL независимо от платформы.
// limit == nil — без ограничения, *limit == 0 — пустой список; offset == 0 — без пропуска.
type ClientRepository interface {
	Create(ctx context.Context, client Client) error
	Delete(ctx context.Context, clientID string) (bool, error)
	FindByName(ctx context.Context, name string, surname string) ([]Client, error)
	List(ctx context.Context, limit *int64, offset int64) ([]Client, error)
	UpdateAddress(ctx context.Context, clientID string, address Address) (Address, bool, error)
}
