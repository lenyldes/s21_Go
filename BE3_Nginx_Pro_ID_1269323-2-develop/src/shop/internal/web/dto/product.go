package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Указатели отличают отсутствующие price и available_stock от допустимого нуля.
type CreateProductRequest struct {
	Name           string   `json:"name" binding:"required" example:"Чайник"`
	Category       string   `json:"category" binding:"required" example:"Кухонная техника"`
	Price          *float64 `json:"price" binding:"required" minimum:"0" example:"2499.99"`
	AvailableStock *int     `json:"available_stock" binding:"required" minimum:"0" maximum:"2147483647" example:"10"`
	LastUpdateDate string   `json:"last_update_date" binding:"required" format:"date-time" example:"2026-10-15T10:00:00Z"`
	SupplierID     string   `json:"supplier_id" binding:"required" format:"uuid" example:"22222222-2222-4222-8222-222222222222"`
	ImageID        *string  `json:"image_id" format:"uuid" example:"33333333-3333-4333-8333-333333333333"`
}

type DecreaseProductStockRequest struct {
	Quantity int `json:"quantity" binding:"required" minimum:"1" maximum:"2147483647" example:"3"`
}

type ListProductsRequest struct {
	Limit  *int64
	Offset int64
}

type ProductResponse struct {
	ID             string  `json:"id" format:"uuid" example:"11111111-1111-4111-8111-111111111111"`
	Name           string  `json:"name" example:"Чайник"`
	Category       string  `json:"category" example:"Кухонная техника"`
	Price          float64 `json:"price" minimum:"0" example:"2499.99"`
	AvailableStock int     `json:"available_stock" minimum:"0" maximum:"2147483647" example:"10"`
	LastUpdateDate string  `json:"last_update_date" format:"date-time" example:"2026-10-15T10:00:00Z"`
	SupplierID     string  `json:"supplier_id" format:"uuid" example:"22222222-2222-4222-8222-222222222222"`
	ImageID        *string `json:"image_id" format:"uuid" example:"33333333-3333-4333-8333-333333333333"`
}

func (dto CreateProductRequest) Validate() error {
	if strings.TrimSpace(dto.Name) == "" {
		return fmt.Errorf("поле name обязательно")
	}
	if strings.TrimSpace(dto.Category) == "" {
		return fmt.Errorf("поле category обязательно")
	}
	if dto.Price == nil {
		return fmt.Errorf("поле price обязательно")
	}
	if *dto.Price < 0 {
		return fmt.Errorf("поле price должно быть неотрицательным числом")
	}
	if dto.AvailableStock == nil {
		return fmt.Errorf("поле available_stock обязательно")
	}
	// Верхняя граница соответствует типу INTEGER в PostgreSQL.
	if *dto.AvailableStock < 0 || *dto.AvailableStock > 2147483647 {
		return fmt.Errorf("поле available_stock должно быть целым числом от 0 до 2147483647")
	}

	_, err := uuid.Parse(dto.SupplierID)
	if err != nil {
		return fmt.Errorf("поле supplier_id должно содержать корректный UUID")
	}
	if dto.ImageID != nil {
		_, err = uuid.Parse(*dto.ImageID)
		if err != nil {
			return fmt.Errorf("поле image_id должно содержать корректный UUID")
		}
	}

	_, err = time.Parse(time.RFC3339, dto.LastUpdateDate)
	if err != nil {
		return fmt.Errorf("поле last_update_date должно содержать корректную дату и время в формате RFC3339")
	}

	return nil
}

func (dto DecreaseProductStockRequest) Validate() error {
	// Верхняя граница соответствует типу INTEGER в PostgreSQL.
	if dto.Quantity < 1 || dto.Quantity > 2147483647 {
		return fmt.Errorf("поле quantity должно быть целым числом от 1 до 2147483647")
	}

	return nil
}

func (dto ListProductsRequest) Validate() error {
	if dto.Limit != nil && *dto.Limit < 0 {
		return fmt.Errorf("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/products?limit=10&offset=20")
	}
	if dto.Offset < 0 {
		return fmt.Errorf("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/products?limit=10&offset=20")
	}

	return nil
}
