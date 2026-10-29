package mapper

import (
	"shop/internal/domain"
	"shop/internal/web/dto"

	"github.com/google/uuid"
)

func CreateAddressToDomain(input dto.CreateAddressRequest) domain.Address {
	return domain.Address{
		ID:      uuid.NewString(),
		Country: input.Country,
		City:    input.City,
		Street:  input.Street,
	}
}

func ReplaceAddressToDomain(addressID string, input dto.ReplaceAddressRequest) domain.Address {
	return domain.Address{
		ID:      addressID,
		Country: input.Country,
		City:    input.City,
		Street:  input.Street,
	}
}

func AddressToWeb(address domain.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:      address.ID,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}
