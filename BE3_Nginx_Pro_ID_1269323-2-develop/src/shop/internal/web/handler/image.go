package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"shop/internal/domain"
	"shop/internal/web/dto"
	"shop/internal/web/mapper"

	"github.com/google/uuid"
)

type ImageHandler struct {
	repository domain.ImageRepository
}

func NewImageHandler(repository domain.ImageRepository) *ImageHandler {
	return &ImageHandler{repository: repository}
}

// @Summary Create product image
// @Description Принимает непустое тело с исходными байтами файла без JSON, base64 и multipart. Формат и содержимое файла не проверяются. UUID изображения генерирует сервер. Изображение создаётся и привязывается к товару в одной транзакции; при отсутствии товара или занятой связи новая запись не сохраняется.
// @Tags Images
// @Security BearerAuth
// @Accept application/octet-stream
// @Produce json
// @Param id path string true "Product ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Param image body string true "Raw image bytes"
// @Success 201 {object} dto.ImageResponse "Image created"
// @Failure 400 {object} dto.ErrorResponse "Invalid product ID, empty body or body read error"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Product not found"
// @Failure 409 {object} dto.ErrorResponse "Product already has an image"
// @Failure 500 {object} dto.ErrorResponse "Failed to create image"
// @Router /products/{id}/image [post]
func (h *ImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Тело содержит исходные байты файла, поэтому JSON-декодирование не требуется.
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, errors.New("ошибка чтения данных изображения"), http.StatusBadRequest)
		return
	}

	request := dto.CreateImageRequest{
		ProductID: r.PathValue("id"),
		Image:     imageData,
	}
	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	image := mapper.ImageToDomain(request)
	productFound, created, err := h.repository.Create(r.Context(), image, request.ProductID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !productFound {
		writeError(w, errors.New("товар не найден"), http.StatusNotFound)
		return
	}
	if !created {
		writeError(w, errors.New("у товара уже есть изображение"), http.StatusConflict)
		return
	}

	response := mapper.ImageToWeb(image)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при создании изображения: %v", err)
	}
}

// @Summary Update image
// @Description Полностью заменяет содержимое изображения исходными байтами непустого тела без JSON, base64 и multipart. Формат и содержимое файла не проверяются. UUID сохраняется; изменение действует для всех товаров, связанных с этим изображением. Изображение можно изменить и без привязки к товару.
// @Tags Images
// @Security BearerAuth
// @Accept application/octet-stream
// @Produce json
// @Param id path string true "Image ID" Format(uuid) example(22222222-2222-4222-8222-222222222222)
// @Param image body string true "Raw image bytes"
// @Success 204 "Image updated"
// @Failure 400 {object} dto.ErrorResponse "Invalid image ID, empty body or body read error"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Image not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to update image"
// @Router /images/{id} [put]
func (h *ImageHandler) Update(w http.ResponseWriter, r *http.Request) {
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, errors.New("ошибка чтения данных изображения"), http.StatusBadRequest)
		return
	}

	request := dto.UpdateImageRequest{
		ImageID: r.PathValue("id"),
		Image:   imageData,
	}
	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	image := mapper.UpdateImageToDomain(request)
	updated, err := h.repository.Update(r.Context(), image)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !updated {
		writeError(w, errors.New("изображение не найдено"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete image
// @Description Удаляет изображение и обнуляет ссылки image_id всех связанных товаров. Сами товары и остальные их поля сохраняются. Можно удалить изображение без привязки к товару.
// @Tags Images
// @Security BearerAuth
// @Produce json
// @Param id path string true "Image ID" Format(uuid) example(22222222-2222-4222-8222-222222222222)
// @Success 204 "Image deleted"
// @Failure 400 {object} dto.ErrorResponse "Invalid image ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Image not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete image"
// @Router /images/{id} [delete]
func (h *ImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	imageID := r.PathValue("id")
	_, err := uuid.Parse(imageID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	deleted, err := h.repository.Delete(r.Context(), imageID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !deleted {
		writeError(w, errors.New("изображение не найдено"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get image by ID
// @Description Возвращает исходные байты файла независимо от наличия связанных товаров. Заголовок Content-Disposition предлагает скачать файл с именем <UUID изображения>.bin. Ошибки возвращаются в JSON без Content-Disposition.
// @Tags Images
// @Security BearerAuth
// @Produce application/octet-stream,json
// @Param id path string true "Image ID" Format(uuid) example(22222222-2222-4222-8222-222222222222)
// @Success 200 {file} file "Image file"
// @Header 200 {string} Content-Disposition "Attachment with the image UUID and .bin extension"
// @Failure 400 {object} dto.ErrorResponse "Invalid image ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Image not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve image"
// @Router /images/{id} [get]
func (h *ImageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	imageID := r.PathValue("id")
	_, err := uuid.Parse(imageID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	image, found, err := h.repository.GetByID(r.Context(), imageID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !found {
		writeError(w, errors.New("изображение не найдено"), http.StatusNotFound)
		return
	}

	response := mapper.ImageFileToWeb(image)
	writeImageFile(w, response)
}

// @Summary Get product image
// @Description Возвращает исходные байты изображения, связанного с товаром. Заголовок Content-Disposition предлагает скачать файл с именем <UUID изображения>.bin; в имени используется UUID изображения, а не товара. Если товар отсутствует или у него нет изображения, возвращается 404. Ошибки возвращаются в JSON без Content-Disposition.
// @Tags Images
// @Security BearerAuth
// @Produce application/octet-stream,json
// @Param id path string true "Product ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Success 200 {file} file "Product image file"
// @Header 200 {string} Content-Disposition "Attachment with the image UUID and .bin extension"
// @Failure 400 {object} dto.ErrorResponse "Invalid product ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Product not found or product has no image"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve product image"
// @Router /products/{id}/image [get]
func (h *ImageHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	_, err := uuid.Parse(productID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	image, productFound, imageFound, err := h.repository.GetByProductID(r.Context(), productID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !productFound {
		writeError(w, errors.New("товар не найден"), http.StatusNotFound)
		return
	}
	if !imageFound {
		writeError(w, errors.New("у товара нет изображения"), http.StatusNotFound)
		return
	}

	response := mapper.ImageFileToWeb(image)
	writeImageFile(w, response)
}

func writeImageFile(w http.ResponseWriter, response dto.ImageFileResponse) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+response.ImageID+".bin\"")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(response.Image)
	if err != nil {
		// Заголовки уже отправлены, поэтому ошибку можно только записать в журнал.
		log.Printf("ошибка отправки файла изображения: %v", err)
	}
}
