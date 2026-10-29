package mapper

import (
	"shop/internal/domain"
	"shop/internal/repository/dao"
)

func ClientToDAO(client domain.Client) dao.Client {
	return dao.Client{
		ID:               client.ID,
		ClientName:       client.ClientName,
		ClientSurname:    client.ClientSurname,
		Birthday:         client.Birthday,
		Gender:           client.Gender,
		RegistrationDate: client.RegistrationDate,
		AddressID:        client.AddressID,
	}
}

func ClientToDomain(client dao.Client) domain.Client {
	return domain.Client{
		ID:               client.ID,
		ClientName:       client.ClientName,
		ClientSurname:    client.ClientSurname,
		Birthday:         client.Birthday,
		Gender:           client.Gender,
		RegistrationDate: client.RegistrationDate,
		AddressID:        client.AddressID,
	}
}
