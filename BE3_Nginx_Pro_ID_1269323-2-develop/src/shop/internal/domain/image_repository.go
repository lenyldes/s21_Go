package domain

import "context"

// В методе Create при отсутствии ошибки productFound == false означает, что товар отсутствует.
// productFound == true и created == false означают, что у товара уже есть изображение.
// В методе Update updated == false при отсутствии ошибки означает, что изображение отсутствует.
// В методе Delete deleted == false при отсутствии ошибки означает, что изображение отсутствует.
// В методе GetByID found == false при отсутствии ошибки означает, что изображение отсутствует.
// В методе GetByProductID при отсутствии ошибки productFound == false означает, что товар отсутствует.
// productFound == true и imageFound == false означают, что у товара нет изображения.
type ImageRepository interface {
	Create(ctx context.Context, image Image, productID string) (productFound bool, created bool, err error)
	Update(ctx context.Context, image Image) (updated bool, err error)
	Delete(ctx context.Context, imageID string) (deleted bool, err error)
	GetByID(ctx context.Context, imageID string) (image Image, found bool, err error)
	GetByProductID(ctx context.Context, productID string) (image Image, productFound bool, imageFound bool, err error)
}
