package mapper

import (
	"shop/internal/domain"
	"shop/internal/repository/dao"
)

func ImageToDAO(image domain.Image) dao.Image {
	return dao.Image{
		ID:    image.ID,
		Image: image.Image,
	}
}

func ImageToDomain(image dao.Image) domain.Image {
	return domain.Image{
		ID:    image.ID,
		Image: image.Image,
	}
}
