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

type ImageRepository struct {
	pool *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) domain.ImageRepository {
	return &ImageRepository{pool: pool}
}

func (r *ImageRepository) Create(ctx context.Context, image domain.Image, productID string) (productFound bool, created bool, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, false, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	// При раннем выходе отменяет изменения; после Commit транзакция уже закрыта.
	defer tx.Rollback(ctx)

	imageDAO := mapper.ImageToDAO(image)
	_, err = tx.Exec(ctx, "INSERT INTO images (id, image) VALUES ($1, $2)", imageDAO.ID, imageDAO.Image)
	if err != nil {
		return false, false, fmt.Errorf("ошибка при сохранении изображения: %w", err)
	}

	// Условие в UPDATE не позволяет заменить уже существующую связь с изображением.
	result, err := tx.Exec(ctx, "UPDATE product SET image_id = $1 WHERE id = $2 AND image_id IS NULL", imageDAO.ID, productID)
	if err != nil {
		return false, false, fmt.Errorf("ошибка при обновлении ссылки на изображение: %w", err)
	}
	if result.RowsAffected() == 0 {
		err = tx.QueryRow(ctx, "SELECT EXISTS (SELECT id FROM product WHERE id = $1)", productID).Scan(&productFound)
		if err != nil {
			return false, false, fmt.Errorf("ошибка при проверке наличия товара: %w", err)
		}
		return productFound, false, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, false, fmt.Errorf("ошибка завершения транзакции: %w", err)
	}

	return true, true, nil
}

func (r *ImageRepository) Update(ctx context.Context, image domain.Image) (updated bool, err error) {
	imageDAO := mapper.ImageToDAO(image)
	query := "UPDATE images SET image = $1 WHERE id = $2"
	result, err := r.pool.Exec(ctx, query, imageDAO.Image, imageDAO.ID)
	if err != nil {
		return false, fmt.Errorf("ошибка при изменении изображения: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (r *ImageRepository) Delete(ctx context.Context, imageID string) (deleted bool, err error) {
	query := "DELETE FROM images WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, imageID)
	if err != nil {
		return false, fmt.Errorf("ошибка при удалении изображения: %w", err)
	}

	return result.RowsAffected() > 0, nil
}

func (r *ImageRepository) GetByID(ctx context.Context, imageID string) (image domain.Image, found bool, err error) {
	query := "SELECT id, image FROM images WHERE id = $1"
	var imageDAO dao.Image
	err = r.pool.QueryRow(ctx, query, imageID).Scan(&imageDAO.ID, &imageDAO.Image)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Image{}, false, nil
	}
	if err != nil {
		return domain.Image{}, false, fmt.Errorf("ошибка при получении изображения: %w", err)
	}

	return mapper.ImageToDomain(imageDAO), true, nil
}

func (r *ImageRepository) GetByProductID(ctx context.Context, productID string) (image domain.Image, productFound bool, imageFound bool, err error) {
	query := `SELECT images.id, images.image
		FROM images
		JOIN product ON product.image_id = images.id
		WHERE product.id = $1`
	var imageDAO dao.Image
	err = r.pool.QueryRow(ctx, query, productID).Scan(&imageDAO.ID, &imageDAO.Image)
	if errors.Is(err, pgx.ErrNoRows) {
		// Отсутствие строки JOIN может означать как отсутствие товара, так и пустую ссылку.
		err = r.pool.QueryRow(ctx, "SELECT EXISTS (SELECT id FROM product WHERE id = $1)", productID).Scan(&productFound)
		if err != nil {
			return domain.Image{}, false, false, fmt.Errorf("ошибка при проверке наличия товара: %w", err)
		}
		return domain.Image{}, productFound, false, nil
	}
	if err != nil {
		return domain.Image{}, false, false, fmt.Errorf("ошибка при получении изображения товара: %w", err)
	}

	return mapper.ImageToDomain(imageDAO), true, true, nil
}
