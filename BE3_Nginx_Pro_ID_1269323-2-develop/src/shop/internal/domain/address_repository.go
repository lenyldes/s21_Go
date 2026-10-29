package domain

import "context"

// В методе Update updated == false при отсутствии ошибки означает, что адрес отсутствует.
// В методе Delete при отсутствии ошибки deleted == false и found == true означают,
// что адрес используется; оба значения false означают, что адрес отсутствует.
// В методе List int64 соответствует диапазону LIMIT/OFFSET PostgreSQL независимо от платформы.
// limit == nil — без ограничения, *limit == 0 — пустой список; offset == 0 — без пропуска.
type AddressRepository interface {
	Create(ctx context.Context, address Address) error
	List(ctx context.Context, limit *int64, offset int64) ([]Address, error)
	Update(ctx context.Context, address Address) (updated bool, err error)
	Delete(ctx context.Context, addressID string) (deleted bool, found bool, err error)
}
