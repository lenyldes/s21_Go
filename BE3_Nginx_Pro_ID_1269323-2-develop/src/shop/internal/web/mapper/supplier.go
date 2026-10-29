package mapper

import (
	"shop/internal/domain"
	"shop/internal/web/dto"

	"github.com/google/uuid"
)

func SupplierToDomain(input dto.CreateSupplierRequest) domain.Supplier {
	return domain.Supplier{
		ID:          uuid.NewString(),
		Name:        input.Name,
		AddressID:   input.AddressID,
		PhoneNumber: input.PhoneNumber,
	}
}

func SupplierToWeb(supplier domain.Supplier) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:          supplier.ID,
		Name:        supplier.Name,
		AddressID:   supplier.AddressID,
		PhoneNumber: supplier.PhoneNumber,
	}
}
