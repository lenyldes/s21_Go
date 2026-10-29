package domain

import "context"

// В методе DecreaseStock при отсутствии ошибки found означает наличие товара, decreased — успешное уменьшение остатка.
// found == false — товар отсутствует; found == true и decreased == false — недостаточный остаток.
// В методе GetByID found == false при отсутствии ошибки означает, что товар отсутствует.
// В методе List int64 соответствует диапазону LIMIT/OFFSET PostgreSQL независимо от платформы.
// limit == nil — без ограничения, *limit == 0 — пустой список; offset == 0 — без пропуска.
type ProductRepository interface {
	Create(ctx context.Context, product Product) error
	Delete(ctx context.Context, productID string) (bool, error)
	DecreaseStock(ctx context.Context, productID string, quantity int) (found bool, decreased bool, err error)
	GetByID(ctx context.Context, productID string) (product Product, found bool, err error)
	List(ctx context.Context, limit *int64, offset int64) ([]Product, error)
}
