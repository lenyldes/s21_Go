package repository

import (
	"context"
	"fmt"

	"shop/internal/domain"
	"shop/internal/repository/dao"
	"shop/internal/repository/mapper"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) domain.ClientRepository {
	return &ClientRepository{pool: pool}
}

func (r *ClientRepository) Create(ctx context.Context, client domain.Client) error {
	query := `INSERT INTO client (
		id, client_name, client_surname, birthday,
		gender, registration_date, address_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	clientDAO := mapper.ClientToDAO(client)
	_, err := r.pool.Exec(ctx, query,
		clientDAO.ID,
		clientDAO.ClientName,
		clientDAO.ClientSurname,
		clientDAO.Birthday,
		clientDAO.Gender,
		clientDAO.RegistrationDate,
		clientDAO.AddressID,
	)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении клиента: %w", err)
	}

	return nil
}

func (r *ClientRepository) Delete(ctx context.Context, clientID string) (bool, error) {
	query := "DELETE FROM client WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, clientID)
	if err != nil {
		return false, fmt.Errorf("ошибка при удалении клиента: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (r *ClientRepository) FindByName(ctx context.Context, name string, surname string) ([]domain.Client, error) {
	query := `SELECT id, client_name, client_surname, birthday,
		gender, registration_date, address_id
		FROM client WHERE client_name = $1 AND client_surname = $2`

	rows, err := r.pool.Query(ctx, query, name, surname)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске клиентов: %w", err)
	}
	defer rows.Close()

	clients := make([]domain.Client, 0)
	for rows.Next() {
		var clientDAO dao.Client
		err = rows.Scan(
			&clientDAO.ID,
			&clientDAO.ClientName,
			&clientDAO.ClientSurname,
			&clientDAO.Birthday,
			&clientDAO.Gender,
			&clientDAO.RegistrationDate,
			&clientDAO.AddressID,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении клиента: %w", err)
		}

		clients = append(clients, mapper.ClientToDomain(clientDAO))
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении списка клиентов: %w", err)
	}

	return clients, nil
}

func (r *ClientRepository) List(ctx context.Context, limit *int64, offset int64) ([]domain.Client, error) {
	query := `SELECT id, client_name, client_surname, birthday,
		gender, registration_date, address_id
		FROM client ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка клиентов: %w", err)
	}
	defer rows.Close()

	clients := make([]domain.Client, 0)
	for rows.Next() {
		var clientDAO dao.Client
		err = rows.Scan(
			&clientDAO.ID,
			&clientDAO.ClientName,
			&clientDAO.ClientSurname,
			&clientDAO.Birthday,
			&clientDAO.Gender,
			&clientDAO.RegistrationDate,
			&clientDAO.AddressID,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении клиента: %w", err)
		}

		clients = append(clients, mapper.ClientToDomain(clientDAO))
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении списка клиентов: %w", err)
	}

	return clients, nil
}

func (r *ClientRepository) UpdateAddress(ctx context.Context, clientID string, address domain.Address) (domain.Address, bool, error) {
	query := "UPDATE client SET address_id = $1 WHERE id = $2"
	currentQuery := `SELECT a.id, a.country, a.city, a.street
		FROM client e JOIN address a ON a.id = e.address_id WHERE e.id = $1`
	savedAddress, updated, err := updateAddress(ctx, r.pool, address, query, currentQuery, clientID)
	if err != nil {
		return domain.Address{}, false, fmt.Errorf("ошибка при изменении адреса клиента: %w", err)
	}

	return savedAddress, updated, nil
}
