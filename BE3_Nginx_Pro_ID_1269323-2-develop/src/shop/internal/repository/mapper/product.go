package mapper

import (
	"shop/internal/domain"
	"shop/internal/repository/dao"
)

func ProductToDAO(product domain.Product) dao.Product {
	return dao.Product{
		ID:             product.ID,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		LastUpdateDate: product.LastUpdateDate,
		SupplierID:     product.SupplierID,
		ImageID:        product.ImageID,
	}
}

func ProductToDomain(product dao.Product) domain.Product {
	return domain.Product{
		ID:             product.ID,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		LastUpdateDate: product.LastUpdateDate,
		SupplierID:     product.SupplierID,
		ImageID:        product.ImageID,
	}
}
