package mapper

import (
	"fmt"
	"time"

	"shop/internal/domain"
	"shop/internal/web/dto"

	"github.com/google/uuid"
)

func ClientToDomain(input dto.CreateClientRequest) (domain.Client, error) {
	birthday, err := time.Parse(time.DateOnly, input.Birthday)
	if err != nil {
		return domain.Client{}, fmt.Errorf("не удалось преобразовать birthday: %w", err)
	}
	registrationDate, err := time.Parse(time.RFC3339, input.RegistrationDate)
	if err != nil {
		return domain.Client{}, fmt.Errorf("не удалось преобразовать registration_date: %w", err)
	}

	return domain.Client{
		ID:               uuid.NewString(),
		ClientName:       input.ClientName,
		ClientSurname:    input.ClientSurname,
		Birthday:         birthday,
		Gender:           input.Gender,
		RegistrationDate: registrationDate,
		AddressID:        input.AddressID,
	}, nil
}

func ClientToWeb(client domain.Client) dto.ClientResponse {
	return dto.ClientResponse{
		ID:               client.ID,
		ClientName:       client.ClientName,
		ClientSurname:    client.ClientSurname,
		Birthday:         client.Birthday.Format(time.DateOnly),
		Gender:           client.Gender,
		RegistrationDate: client.RegistrationDate.Format(time.RFC3339),
		AddressID:        client.AddressID,
	}
}
