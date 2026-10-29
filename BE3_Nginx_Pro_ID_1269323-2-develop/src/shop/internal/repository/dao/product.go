package dao

import "time"

type Product struct {
	ID             string
	Name           string
	Category       string
	Price          float64
	AvailableStock int
	LastUpdateDate time.Time
	SupplierID     string
	ImageID        *string // nil, если изображение ещё не добавлено.
}
