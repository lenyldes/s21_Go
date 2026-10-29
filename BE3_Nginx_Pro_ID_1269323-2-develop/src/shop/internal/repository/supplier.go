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

type SupplierRepository struct {
	pool *pgxpool.Pool
}

func NewSupplierRepository(pool *pgxpool.Pool) domain.SupplierRepository {
	return &SupplierRepository{pool: pool}
}

func (r *SupplierRepository) Create(ctx context.Context, supplier domain.Supplier) error {
	query := "INSERT INTO supplier (id, name, address_id, phone_number) VALUES ($1, $2, $3, $4)"

	supplierDAO := mapper.SupplierToDAO(supplier)
	_, err := r.pool.Exec(ctx, query,
		supplierDAO.ID,
		supplierDAO.Name,
		supplierDAO.AddressID,
		supplierDAO.PhoneNumber,
	)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении поставщика: %w", err)
	}

	return nil
}

func (r *SupplierRepository) Delete(ctx context.Context, supplierID string) (bool, error) {
	query := "DELETE FROM supplier WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, supplierID)
	if err != nil {
		return false, fmt.Errorf("ошибка при удалении поставщика: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (r *SupplierRepository) UpdateAddress(ctx context.Context, supplierID string, address domain.Address) (domain.Address, bool, error) {
	query := "UPDATE supplier SET address_id = $1 WHERE id = $2"
	currentQuery := `SELECT a.id, a.country, a.city, a.street
		FROM supplier e JOIN address a ON a.id = e.address_id WHERE e.id = $1`
	savedAddress, updated, err := updateAddress(ctx, r.pool, address, query, currentQuery, supplierID)
	if err != nil {
		return domain.Address{}, false, fmt.Errorf("ошибка при изменении адреса поставщика: %w", err)
	}

	return savedAddress, updated, nil
}

func (r *SupplierRepository) GetByID(ctx context.Context, supplierID string) (domain.Supplier, bool, error) {
	query := "SELECT id, name, address_id, phone_number FROM supplier WHERE id = $1"

	var supplierDAO dao.Supplier
	err := r.pool.QueryRow(ctx, query, supplierID).Scan(
		&supplierDAO.ID,
		&supplierDAO.Name,
		&supplierDAO.AddressID,
		&supplierDAO.PhoneNumber,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Supplier{}, false, nil
	}
	if err != nil {
		return domain.Supplier{}, false, fmt.Errorf("ошибка при получении поставщика: %w", err)
	}

	return mapper.SupplierToDomain(supplierDAO), true, nil
}

func (r *SupplierRepository) List(ctx context.Context, limit *int64, offset int64) ([]domain.Supplier, error) {
	query := "SELECT id, name, address_id, phone_number FROM supplier ORDER BY id LIMIT $1 OFFSET $2"

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка поставщиков: %w", err)
	}
	defer rows.Close()

	suppliers := make([]domain.Supplier, 0)
	for rows.Next() {
		var supplierDAO dao.Supplier
		err = rows.Scan(
			&supplierDAO.ID,
			&supplierDAO.Name,
			&supplierDAO.AddressID,
			&supplierDAO.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении поставщика: %w", err)
		}

		suppliers = append(suppliers, mapper.SupplierToDomain(supplierDAO))
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении списка поставщиков: %w", err)
	}

	return suppliers, nil
}
