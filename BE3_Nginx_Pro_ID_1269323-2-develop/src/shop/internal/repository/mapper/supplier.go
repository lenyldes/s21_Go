package mapper

import (
	"shop/internal/domain"
	"shop/internal/repository/dao"
)

func SupplierToDAO(supplier domain.Supplier) dao.Supplier {
	return dao.Supplier{
		ID:          supplier.ID,
		Name:        supplier.Name,
		AddressID:   supplier.AddressID,
		PhoneNumber: supplier.PhoneNumber,
	}
}

func SupplierToDomain(supplier dao.Supplier) domain.Supplier {
	return domain.Supplier{
		ID:          supplier.ID,
		Name:        supplier.Name,
		AddressID:   supplier.AddressID,
		PhoneNumber: supplier.PhoneNumber,
	}
}
