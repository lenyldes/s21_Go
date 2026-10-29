package dto

import (
	"fmt"

	"github.com/google/uuid"
)

// CreateImageRequest заполняется из ID товара в пути и байтов тела запроса.
type CreateImageRequest struct {
	ProductID string
	Image     []byte
}

// UpdateImageRequest заполняется из ID изображения в пути и новых байтов тела запроса.
type UpdateImageRequest struct {
	ImageID string
	Image   []byte
}

type ImageResponse struct {
	ID string `json:"id" format:"uuid" example:"22222222-2222-4222-8222-222222222222"`
}

// ImageFileResponse передаётся как файл: поле Image записывается в тело без JSON-кодирования.
type ImageFileResponse struct {
	ImageID string
	Image   []byte
}

func (dto CreateImageRequest) Validate() error {
	_, err := uuid.Parse(dto.ProductID)
	if err != nil {
		return fmt.Errorf("поле id должно содержать корректный UUID")
	}
	if len(dto.Image) == 0 {
		return fmt.Errorf("данные изображения не должны быть пустыми")
	}

	return nil
}

func (dto UpdateImageRequest) Validate() error {
	_, err := uuid.Parse(dto.ImageID)
	if err != nil {
		return fmt.Errorf("поле id должно содержать корректный UUID")
	}
	if len(dto.Image) == 0 {
		return fmt.Errorf("данные изображения не должны быть пустыми")
	}

	return nil
}
