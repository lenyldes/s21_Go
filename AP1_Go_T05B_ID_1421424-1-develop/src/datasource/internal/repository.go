package internal

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"tictactoe/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepository struct {
	db *pgxpool.Pool
}

func NewGameRepository(db *pgxpool.Pool) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) GetAvailableGames(ctx context.Context) ([]domain.Game, error) {
	result := make([]domain.Game, 0)

	query := "SELECT uuid, board, phase, player_uuid, player_x, player_o, created_at FROM games WHERE phase = $1"

	rows, err := r.db.Query(ctx, query, domain.GamePhaseReady)
	if err != nil {
		return []domain.Game{}, fmt.Errorf("ошибка при получении списка доступных игр: %w", err)
	}
	defer rows.Close()

	for rows.Next() == true {
		var gameEntity GameEntity
		err = rows.Scan(&gameEntity.UUID, &gameEntity.Board, &gameEntity.Phase, &gameEntity.PlayerUUID, &gameEntity.PlayerX, &gameEntity.PlayerO, &gameEntity.CreatedAt)
		if err != nil {
			return []domain.Game{}, fmt.Errorf("ошибка при получении списка доступных игр: %w", err)
		}
		result = append(result, gameToDomain(gameEntity))
	}
	if rows.Err() != nil {
		return []domain.Game{}, fmt.Errorf("ошибка при получении списка доступных игр: %w", rows.Err())
	}

	return result, nil
}

func (r *GameRepository) GetCompletedGames(ctx context.Context, userUUID string) ([]domain.Game, error) {
	result := make([]domain.Game, 0)

	query := "SELECT uuid, board, phase, player_uuid, player_x, player_o, created_at FROM games WHERE (player_x = $1 OR player_o = $1) AND phase IN ($2, $3)"

	rows, err := r.db.Query(ctx, query, userUUID, domain.GamePhaseWin, domain.GamePhaseDraw)
	if err != nil {
		return []domain.Game{}, fmt.Errorf("ошибка при получении списка завершенных игр: %w", err)
	}
	defer rows.Close()

	for rows.Next() == true {
		var gameEntity GameEntity
		err = rows.Scan(&gameEntity.UUID, &gameEntity.Board, &gameEntity.Phase, &gameEntity.PlayerUUID, &gameEntity.PlayerX, &gameEntity.PlayerO, &gameEntity.CreatedAt)
		if err != nil {
			return []domain.Game{}, fmt.Errorf("ошибка при получении списка завершенных игр: %w", err)
		}
		result = append(result, gameToDomain(gameEntity))
	}
	if rows.Err() != nil {
		return []domain.Game{}, fmt.Errorf("ошибка при получении списка завершенных игр: %w", rows.Err())
	}

	return result, nil
}

func (r *GameRepository) getUserUUIDs(ctx context.Context) ([]string, error) {
	query := "SELECT uuid FROM users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []string{}, fmt.Errorf("ошибка при получении UUID пользователей: %w", err)
	}
	defer rows.Close()

	userUUIDs := make([]string, 0)
	for rows.Next() {
		var userUUID string
		if err := rows.Scan(&userUUID); err != nil {
			return []string{}, fmt.Errorf("ошибка при получении UUID пользователей: %w", err)
		}
		userUUIDs = append(userUUIDs, userUUID)
	}
	if err := rows.Err(); err != nil {
		return []string{}, fmt.Errorf("ошибка при получении UUID пользователей: %w", err)
	}

	return userUUIDs, nil
}

func (r *GameRepository) GetTopPlayers(ctx context.Context, n int) ([]domain.WinStatistics, error) {
	if n < 0 {
		return []domain.WinStatistics{}, fmt.Errorf("количество лучших игроков не может быть отрицательным")
	}

	userUUIDs, err := r.getUserUUIDs(ctx)
	if err != nil {
		return []domain.WinStatistics{}, fmt.Errorf("ошибка при получении списка лучших игроков: %w", err)
	}

	result := make([]domain.WinStatistics, 0)

	for _, userUUID := range userUUIDs {
		games, err := r.GetCompletedGames(ctx, userUUID)
		if err != nil {
			return []domain.WinStatistics{}, fmt.Errorf("ошибка при получении списка лучших игроков: %w", err)
		}

		wins, losses, draws := 0, 0, 0
		for _, game := range games {
			if game.State.Phase == domain.GamePhaseDraw {
				draws++
			} else if game.State.PlayerUUID == userUUID {
				wins++
			} else {
				losses++
			}
		}

		denominator := losses + draws
		if denominator == 0 {
			denominator = 1
		}

		winRatio := float64(wins) / float64(denominator)
		result = append(result, domain.WinStatistics{
			UserUUID: userUUID,
			WinRatio: winRatio,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].WinRatio == result[j].WinRatio {
			return result[i].UserUUID < result[j].UserUUID
		}
		return result[i].WinRatio > result[j].WinRatio
	})

	if n < len(result) {
		result = result[:n]
	}

	return result, nil
}

func (r *GameRepository) Save(ctx context.Context, game domain.Game) error {
	entityGame := gameToEntity(game)

	query := "UPDATE games SET board = $1, phase = $2, player_uuid = $3, player_x = $4, player_o = $5 WHERE uuid = $6"

	tag, err := r.db.Exec(ctx, query, entityGame.Board, entityGame.Phase, entityGame.PlayerUUID, entityGame.PlayerX, entityGame.PlayerO, entityGame.UUID)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении игры: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrGameNotFound
	}

	return nil
}

func (r *GameRepository) Get(ctx context.Context, UUID string) (domain.Game, error) {
	query := "SELECT uuid, board, phase, player_uuid, player_x, player_o, created_at FROM games WHERE uuid = $1"

	var entityGame GameEntity
	err := r.db.QueryRow(ctx, query, UUID).Scan(&entityGame.UUID, &entityGame.Board, &entityGame.Phase, &entityGame.PlayerUUID, &entityGame.PlayerX, &entityGame.PlayerO, &entityGame.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Game{}, domain.ErrGameNotFound
		}
		return domain.Game{}, fmt.Errorf("ошибка запроса: %w", err)
	}

	return gameToDomain(entityGame), nil
}

func (r *GameRepository) Create(ctx context.Context, game domain.Game) (string, error) {
	entityGame := gameToEntity(game)

	query := "INSERT INTO games (uuid, board, phase, player_uuid, player_x, player_o, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, query, entityGame.UUID, entityGame.Board, entityGame.Phase, entityGame.PlayerUUID, entityGame.PlayerX, entityGame.PlayerO, entityGame.CreatedAt)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании игры: %w", err)
	}

	return entityGame.UUID, nil
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (string, error) {
	entityUser := userToEntity(user)

	query := "INSERT INTO users (uuid, login, password) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, query, entityUser.UUID, entityUser.Login, entityUser.Password)
	if err != nil {
		return "", fmt.Errorf("не удалось добавить пользователя в таблицу users: %w", err)
	}

	return entityUser.UUID, nil
}

func (r *UserRepository) GetByUUID(ctx context.Context, UUID string) (domain.User, error) {
	query := "SELECT uuid, login, password FROM users WHERE uuid = $1"

	var entityUser UserEntity

	err := r.db.QueryRow(ctx, query, UUID).Scan(&entityUser.UUID, &entityUser.Login, &entityUser.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("ошибка запроса: %w", err)
	}

	return userToDomain(entityUser), nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	query := "SELECT uuid, login, password FROM users WHERE login = $1"

	var entityUser UserEntity

	err := r.db.QueryRow(ctx, query, login).Scan(&entityUser.UUID, &entityUser.Login, &entityUser.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("ошибка запроса: %w", err)
	}

	return userToDomain(entityUser), nil
}
