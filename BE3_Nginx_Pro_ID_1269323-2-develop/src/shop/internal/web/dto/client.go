package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateClientRequest struct {
	ClientName       string `json:"client_name" binding:"required" example:"Иван"`
	ClientSurname    string `json:"client_surname" binding:"required" example:"Иванов"`
	Birthday         string `json:"birthday" binding:"required" format:"date" example:"1995-05-20"`
	Gender           string `json:"gender" binding:"required" example:"male"`
	RegistrationDate string `json:"registration_date" binding:"required" format:"date-time" example:"2026-10-14T12:00:00Z"`
	AddressID        string `json:"address_id" binding:"required" format:"uuid" example:"22222222-2222-4222-8222-222222222222"`
}

type FindClientsRequest struct {
	ClientName    string
	ClientSurname string
}

type ListClientsRequest struct {
	Limit  *int64
	Offset int64
}

type ClientResponse struct {
	ID               string `json:"id" format:"uuid" example:"11111111-1111-4111-8111-111111111111"`
	ClientName       string `json:"client_name" example:"Иван"`
	ClientSurname    string `json:"client_surname" example:"Иванов"`
	Birthday         string `json:"birthday" format:"date" example:"1995-05-20"`
	Gender           string `json:"gender" example:"male"`
	RegistrationDate string `json:"registration_date" format:"date-time" example:"2026-10-14T12:00:00Z"`
	AddressID        string `json:"address_id" format:"uuid" example:"22222222-2222-4222-8222-222222222222"`
}

func (dto CreateClientRequest) Validate() error {
	if strings.TrimSpace(dto.ClientName) == "" {
		return fmt.Errorf("поле client_name обязательно")
	}
	if strings.TrimSpace(dto.ClientSurname) == "" {
		return fmt.Errorf("поле client_surname обязательно")
	}
	if strings.TrimSpace(dto.Gender) == "" {
		return fmt.Errorf("поле gender обязательно")
	}

	_, err := uuid.Parse(dto.AddressID)
	if err != nil {
		return fmt.Errorf("поле address_id должно содержать корректный UUID")
	}

	_, err = time.Parse(time.DateOnly, dto.Birthday)
	if err != nil {
		return fmt.Errorf("поле birthday должно содержать корректную дату в формате YYYY-MM-DD")
	}
	_, err = time.Parse(time.RFC3339, dto.RegistrationDate)
	if err != nil {
		return fmt.Errorf("поле registration_date должно содержать корректную дату и время в формате RFC3339")
	}

	return nil
}

func (dto FindClientsRequest) Validate() error {
	if strings.TrimSpace(dto.ClientName) == "" {
		return fmt.Errorf("поле client_name обязательно. Пример запроса: /api/v1/clients?client_name=Иван&client_surname=Иванов")
	}
	if strings.TrimSpace(dto.ClientSurname) == "" {
		return fmt.Errorf("поле client_surname обязательно. Пример запроса: /api/v1/clients?client_name=Иван&client_surname=Иванов")
	}

	return nil
}

func (dto ListClientsRequest) Validate() error {
	if dto.Limit != nil && *dto.Limit < 0 {
		return fmt.Errorf("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/clients?limit=10&offset=20")
	}
	if dto.Offset < 0 {
		return fmt.Errorf("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/clients?limit=10&offset=20")
	}

	return nil
}
