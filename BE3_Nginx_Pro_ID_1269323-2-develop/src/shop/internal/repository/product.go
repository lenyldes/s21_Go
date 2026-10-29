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

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) domain.ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, product domain.Product) error {
	query := `INSERT INTO product (
		id, name, category, price,
		available_stock, last_update_date, supplier_id, image_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	productDAO := mapper.ProductToDAO(product)
	_, err := r.pool.Exec(ctx, query,
		productDAO.ID,
		productDAO.Name,
		productDAO.Category,
		productDAO.Price,
		productDAO.AvailableStock,
		productDAO.LastUpdateDate,
		productDAO.SupplierID,
		productDAO.ImageID,
	)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении товара: %w", err)
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, productID string) (bool, error) {
	query := "DELETE FROM product WHERE id = $1"

	result, err := r.pool.Exec(ctx, query, productID)
	if err != nil {
		return false, fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (r *ProductRepository) DecreaseStock(ctx context.Context, productID string, quantity int) (bool, bool, error) {
	var found bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS (SELECT id FROM product WHERE id = $1)", productID).Scan(&found)
	if err != nil {
		return false, false, fmt.Errorf("ошибка при проверке наличия товара: %w", err)
	}
	if !found {
		return false, false, nil
	}

	// Проверка остатка и уменьшение выполняются одним запросом, чтобы избежать списания сверх остатка.
	query := `UPDATE product
		SET available_stock = available_stock - $1
		WHERE id = $2 AND available_stock >= $1`

	result, err := r.pool.Exec(ctx, query, quantity, productID)
	if err != nil {
		return false, false, fmt.Errorf("ошибка при уменьшении количества товара: %w", err)
	}
	return true, result.RowsAffected() > 0, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, productID string) (domain.Product, bool, error) {
	query := `SELECT id, name, category, price,
		available_stock, last_update_date, supplier_id, image_id
		FROM product WHERE id = $1`

	var productDAO dao.Product
	err := r.pool.QueryRow(ctx, query, productID).Scan(
		&productDAO.ID,
		&productDAO.Name,
		&productDAO.Category,
		&productDAO.Price,
		&productDAO.AvailableStock,
		&productDAO.LastUpdateDate,
		&productDAO.SupplierID,
		&productDAO.ImageID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, false, nil
	}
	if err != nil {
		return domain.Product{}, false, fmt.Errorf("ошибка при получении товара: %w", err)
	}

	return mapper.ProductToDomain(productDAO), true, nil
}

func (r *ProductRepository) List(ctx context.Context, limit *int64, offset int64) ([]domain.Product, error) {
	query := `SELECT id, name, category, price,
		available_stock, last_update_date, supplier_id, image_id
		FROM product WHERE available_stock > 0 ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка товаров: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var productDAO dao.Product
		err = rows.Scan(
			&productDAO.ID,
			&productDAO.Name,
			&productDAO.Category,
			&productDAO.Price,
			&productDAO.AvailableStock,
			&productDAO.LastUpdateDate,
			&productDAO.SupplierID,
			&productDAO.ImageID,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении товара: %w", err)
		}

		products = append(products, mapper.ProductToDomain(productDAO))
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении списка товаров: %w", err)
	}

	return products, nil
}
