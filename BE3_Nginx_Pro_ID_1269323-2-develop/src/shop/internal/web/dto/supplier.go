package dto

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type CreateSupplierRequest struct {
	Name        string `json:"name" binding:"required" example:"Бытовая техника"`
	AddressID   string `json:"address_id" binding:"required" format:"uuid" example:"22222222-2222-4222-8222-222222222222"`
	PhoneNumber string `json:"phone_number" binding:"required" example:"+7 (900) 123-45-67"`
}

type ListSuppliersRequest struct {
	Limit  *int64
	Offset int64
}

type SupplierResponse struct {
	ID          string `json:"id" format:"uuid" example:"11111111-1111-4111-8111-111111111111"`
	Name        string `json:"name" example:"Бытовая техника"`
	AddressID   string `json:"address_id" format:"uuid" example:"22222222-2222-4222-8222-222222222222"`
	PhoneNumber string `json:"phone_number" example:"+7 (900) 123-45-67"`
}

func (dto CreateSupplierRequest) Validate() error {
	if strings.TrimSpace(dto.Name) == "" {
		return fmt.Errorf("поле name обязательно")
	}
	if strings.TrimSpace(dto.PhoneNumber) == "" {
		return fmt.Errorf("поле phone_number обязательно")
	}

	_, err := uuid.Parse(dto.AddressID)
	if err != nil {
		return fmt.Errorf("поле address_id должно содержать корректный UUID")
	}

	return nil
}

func (dto ListSuppliersRequest) Validate() error {
	if dto.Limit != nil && *dto.Limit < 0 {
		return fmt.Errorf("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/suppliers?limit=10&offset=20")
	}
	if dto.Offset < 0 {
		return fmt.Errorf("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/suppliers?limit=10&offset=20")
	}

	return nil
}
