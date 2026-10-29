package mapper

import (
	"shop/internal/domain"
	"shop/internal/repository/dao"
)

func AddressToDAO(address domain.Address) dao.Address {
	return dao.Address{
		ID:      address.ID,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}

func AddressToDomain(address dao.Address) domain.Address {
	return domain.Address{
		ID:      address.ID,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}
