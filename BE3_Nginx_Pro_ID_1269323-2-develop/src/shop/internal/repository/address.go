package repository

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/domain"
	"shop/internal/repository/dao"
	"shop/internal/repository/mapper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepository struct {
	pool *pgxpool.Pool
}

func NewAddressRepository(pool *pgxpool.Pool) domain.AddressRepository {
	return &AddressRepository{pool: pool}
}

func (r *AddressRepository) Create(ctx context.Context, address domain.Address) error {
	addressDAO := mapper.AddressToDAO(address)
	query := "INSERT INTO address (id, country, city, street) VALUES ($1, $2, $3, $4)"
	_, err := r.pool.Exec(ctx, query, addressDAO.ID, addressDAO.Country, addressDAO.City, addressDAO.Street)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении адреса: %w", err)
	}

	return nil
}

func (r *AddressRepository) List(ctx context.Context, limit *int64, offset int64) ([]domain.Address, error) {
	query := "SELECT id, country, city, street FROM address ORDER BY id LIMIT $1 OFFSET $2"

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка адресов: %w", err)
	}
	defer rows.Close()

	addresses := make([]domain.Address, 0)
	for rows.Next() {
		var addressDAO dao.Address
		err = rows.Scan(&addressDAO.ID, &addressDAO.Country, &addressDAO.City, &addressDAO.Street)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении адреса: %w", err)
		}

		addresses = append(addresses, mapper.AddressToDomain(addressDAO))
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении списка адресов: %w", err)
	}

	return addresses, nil
}

func (r *AddressRepository) Update(ctx context.Context, address domain.Address) (updated bool, err error) {
	addressDAO := mapper.AddressToDAO(address)
	query := "UPDATE address SET country = $1, city = $2, street = $3 WHERE id = $4"
	result, err := r.pool.Exec(ctx, query, addressDAO.Country, addressDAO.City, addressDAO.Street, addressDAO.ID)
	if err != nil {
		return false, fmt.Errorf("ошибка при изменении адреса: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (r *AddressRepository) Delete(ctx context.Context, addressID string) (deleted bool, found bool, err error) {
	query := `DELETE FROM address WHERE id = $1
		AND NOT EXISTS (SELECT 1 FROM client WHERE address_id = $1)
		AND NOT EXISTS (SELECT 1 FROM supplier WHERE address_id = $1)`
	result, err := r.pool.Exec(ctx, query, addressID)
	if err != nil {
		return false, false, fmt.Errorf("ошибка при удалении адреса: %w", err)
	}
	if result.RowsAffected() > 0 {
		return true, true, nil
	}

	err = r.pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM address WHERE id = $1)", addressID).Scan(&found)
	if err != nil {
		return false, false, fmt.Errorf("ошибка при проверке наличия адреса: %w", err)
	}

	return false, found, nil
}

// updateQuery — фиксированный запрос репозитория: $1 — ID нового адреса, $2 — ID клиента или поставщика.
// currentQuery — фиксированный запрос текущего адреса: $1 — ID клиента или поставщика.
func updateAddress(ctx context.Context, pool *pgxpool.Pool, address domain.Address, updateQuery string, currentQuery string, entityID string) (domain.Address, bool, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return domain.Address{}, false, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	// При раннем выходе отменяет изменения; после Commit транзакция уже закрыта.
	defer tx.Rollback(ctx)

	// Читаем адрес до вставки: для отсутствующей сущности и при совпадении новый адрес не нужен.
	var current dao.Address
	err = tx.QueryRow(ctx, currentQuery, entityID).Scan(&current.ID, &current.Country, &current.City, &current.Street)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Address{}, false, nil
	}
	if err != nil {
		return domain.Address{}, false, fmt.Errorf("ошибка при чтении текущего адреса: %w", err)
	}
	if current.Country == address.Country && current.City == address.City && current.Street == address.Street {
		return mapper.AddressToDomain(current), true, nil
	}

	addressDAO := mapper.AddressToDAO(address)
	query := "INSERT INTO address (id, country, city, street) VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(ctx, query, addressDAO.ID, addressDAO.Country, addressDAO.City, addressDAO.Street)
	if err != nil {
		return domain.Address{}, false, fmt.Errorf("ошибка при сохранении нового адреса: %w", err)
	}

	result, err := tx.Exec(ctx, updateQuery, addressDAO.ID, entityID)
	if err != nil {
		return domain.Address{}, false, fmt.Errorf("ошибка при обновлении ссылки на адрес: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.Address{}, false, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Address{}, false, fmt.Errorf("ошибка завершения транзакции: %w", err)
	}

	return address, true, nil
}
