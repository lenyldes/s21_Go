package domain

import "context"

// В методе GetByID found == false при отсутствии ошибки означает, что поставщик отсутствует.
// В методе List int64 соответствует диапазону LIMIT/OFFSET PostgreSQL независимо от платформы.
// limit == nil — без ограничения, *limit == 0 — пустой список; offset == 0 — без пропуска.
type SupplierRepository interface {
	Create(ctx context.Context, supplier Supplier) error
	Delete(ctx context.Context, supplierID string) (bool, error)
	UpdateAddress(ctx context.Context, supplierID string, address Address) (Address, bool, error)
	GetByID(ctx context.Context, supplierID string) (supplier Supplier, found bool, err error)
	List(ctx context.Context, limit *int64, offset int64) ([]Supplier, error)
}
