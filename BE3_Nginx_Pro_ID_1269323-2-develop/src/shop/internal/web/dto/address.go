package dto

import (
	"fmt"
	"strings"
)

type CreateAddressRequest struct {
	Country string `json:"country" binding:"required" example:"Russia"`
	City    string `json:"city" binding:"required" example:"Moscow"`
	Street  string `json:"street" binding:"required" example:"Tverskaya Street, 1"`
}

type ReplaceAddressRequest struct {
	Country string `json:"country" binding:"required" example:"Russia"`
	City    string `json:"city" binding:"required" example:"Moscow"`
	Street  string `json:"street" binding:"required" example:"Tverskaya Street, 1"`
}

type ListAddressesRequest struct {
	Limit  *int64
	Offset int64
}

type AddressResponse struct {
	ID      string `json:"id" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Country string `json:"country" example:"Russia"`
	City    string `json:"city" example:"Moscow"`
	Street  string `json:"street" example:"Tverskaya Street, 1"`
}

func (dto CreateAddressRequest) Validate() error {
	return (ReplaceAddressRequest{Country: dto.Country, City: dto.City, Street: dto.Street}).Validate()
}

func (dto ReplaceAddressRequest) Validate() error {
	if strings.TrimSpace(dto.Country) == "" {
		return fmt.Errorf("поле country обязательно")
	}
	if strings.TrimSpace(dto.City) == "" {
		return fmt.Errorf("поле city обязательно")
	}
	if strings.TrimSpace(dto.Street) == "" {
		return fmt.Errorf("поле street обязательно")
	}

	return nil
}

func (dto ListAddressesRequest) Validate() error {
	if dto.Limit != nil && *dto.Limit < 0 {
		return fmt.Errorf("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/addresses?limit=10&offset=20")
	}
	if dto.Offset < 0 {
		return fmt.Errorf("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/addresses?limit=10&offset=20")
	}

	return nil
}
