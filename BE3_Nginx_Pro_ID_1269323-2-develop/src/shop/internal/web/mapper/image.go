package mapper

import (
	"shop/internal/domain"
	"shop/internal/web/dto"

	"github.com/google/uuid"
)

func ImageToDomain(input dto.CreateImageRequest) domain.Image {
	return domain.Image{
		ID:    uuid.NewString(),
		Image: input.Image,
	}
}

func ImageToWeb(image domain.Image) dto.ImageResponse {
	return dto.ImageResponse{
		ID: image.ID,
	}
}

func UpdateImageToDomain(input dto.UpdateImageRequest) domain.Image {
	return domain.Image{
		ID:    input.ImageID,
		Image: input.Image,
	}
}

func ImageFileToWeb(image domain.Image) dto.ImageFileResponse {
	return dto.ImageFileResponse{
		ImageID: image.ID,
		Image:   image.Image,
	}
}
