package data

import (
	"auth/internal/auth"
	"auth/internal/shared"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *repo {
	return &repo{db: pool}
}

func InitPostgresPool(config *shared.Config) *pgxpool.Pool {
	fmt.Println(config.PgHost, config.AppPort)
	ctx := context.Background()
	conn_url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.PgUser,
		config.PgPassword,
		config.PgHost,
		config.PgPort,
		config.PgDatabase,
	)
	pool, err := pgxpool.New(ctx, conn_url)
	if err != nil {
		panic(fmt.Sprintf("Unable to create connection pool: %v", err.Error()))
	}
	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("Unable to ping database: %v", err.Error()))
	}
	return pool
}

func (r *repo) SaveUser(ctx context.Context, user *auth.User) error {
	query := `INSERT INTO "user"(id, login, name, surname, phonenumber, passhash, email, role) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(ctx, query, toPgUUID(user.Id), user.Login, user.Name, user.Surname, user.Phonenumber, user.PasswordHash, user.Email, user.Role)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Println("ERR: ", err, pgErr.ColumnName)
		if pgErr.Code == "23505" {
			return auth.ErrLoginTaken
		}
		return fmt.Errorf("Can not execute query: %w", err)
	}

	return nil
}

func (r *repo) GetUserById(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	query := `SELECT * FROM "user" WHERE id = $1`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("Can not execute query: %w", err)
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userRecord])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, auth.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal from database: %w", err)
	}

	return toUser(record), nil
}

func (r *repo) GetUserByLogin(ctx context.Context, login string) (*auth.User, error) {
	query := `SELECT * FROM "user" WHERE login = $1`

	rows, err := r.db.Query(ctx, query, login)
	if err != nil {
		return nil, fmt.Errorf("Can not execute query: %w", err)
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userRecord])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, auth.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal from database: %w", err)
	}

	return toUser(record), nil
}

func (r *repo) ChangeUserPassword(ctx context.Context, id uuid.UUID, passHash []byte) error {
	query := `UPDATE "user" SET passhash = $1 WHERE id = $2`

	tag, err := r.db.Exec(ctx, query, passHash, toPgUUID(id))
	if err != nil {
		return fmt.Errorf("Can not execute query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return auth.ErrUserNotFound
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "user" WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, toPgUUID(id))
	if err != nil {
		return fmt.Errorf("Can not execute query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return auth.ErrUserNotFound
	}

	return nil
}
