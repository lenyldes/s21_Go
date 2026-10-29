package mapper

import (
	"fmt"
	"time"

	"shop/internal/domain"
	"shop/internal/web/dto"

	"github.com/google/uuid"
)

func ProductToDomain(input dto.CreateProductRequest) (domain.Product, error) {
	lastUpdateDate, err := time.Parse(time.RFC3339, input.LastUpdateDate)
	if err != nil {
		return domain.Product{}, fmt.Errorf("не удалось преобразовать last_update_date: %w", err)
	}

	return domain.Product{
		ID:             uuid.NewString(),
		Name:           input.Name,
		Category:       input.Category,
		Price:          *input.Price,
		AvailableStock: *input.AvailableStock,
		LastUpdateDate: lastUpdateDate,
		SupplierID:     input.SupplierID,
		ImageID:        input.ImageID,
	}, nil
}

func ProductToWeb(product domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:             product.ID,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		LastUpdateDate: product.LastUpdateDate.Format(time.RFC3339),
		SupplierID:     product.SupplierID,
		ImageID:        product.ImageID,
	}
}
