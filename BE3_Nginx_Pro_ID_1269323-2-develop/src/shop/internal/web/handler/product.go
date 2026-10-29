package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"shop/internal/domain"
	"shop/internal/web/dto"
	"shop/internal/web/mapper"

	"github.com/google/uuid"
)

type ProductHandler struct {
	repository domain.ProductRepository
}

func NewProductHandler(repository domain.ProductRepository) *ProductHandler {
	return &ProductHandler{repository: repository}
}

// @Summary Create product
// @Description UUID товара генерирует сервер, дату последней закупки передаёт отправитель. Поставщик и указанное изображение должны уже существовать. Поле image_id можно не передавать или задать null; в этом случае товар создаётся без изображения. Название и категория не могут состоять только из пробелов; список категорий не ограничен. Цена и остаток обязательны, но могут быть нулевыми.
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product"
// @Success 201 {object} dto.ProductResponse "Product created"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Failed to create product, including a nonexistent supplier or image"
// @Router /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateProductRequest
	err := decodeJSONBody(r.Body, &request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	product, err := mapper.ProductToDomain(request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = h.repository.Create(r.Context(), product)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := mapper.ProductToWeb(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при создании товара: %v", err)
	}
}

// @Summary Delete product
// @Description Удаление допускается при любом остатке, включая нулевой. Поставщик и изображение сохраняются.
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Success 204 "Product deleted"
// @Failure 400 {object} dto.ErrorResponse "Invalid product ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Product not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete product"
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	_, err := uuid.Parse(productID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	deleted, err := h.repository.Delete(r.Context(), productID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !deleted {
		writeError(w, errors.New("товар не найден"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get product by ID
// @Description Возвращает товар независимо от остатка, включая нулевой. У товара без изображения поле image_id равно null.
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Success 200 {object} dto.ProductResponse "Product"
// @Failure 400 {object} dto.ErrorResponse "Invalid product ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Product not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve product"
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	_, err := uuid.Parse(productID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	product, found, err := h.repository.GetByID(r.Context(), productID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !found {
		writeError(w, errors.New("товар не найден"), http.StatusNotFound)
		return
	}

	response := mapper.ProductToWeb(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при получении товара: %v", err)
	}
}

// @Summary List available products
// @Description Возвращает только товары с available_stock > 0, отсортированные по ID по возрастанию. Параметры пагинации необязательны и могут использоваться отдельно; отбор доступных товаров выполняется до пагинации. Без limit количество результатов не ограничено; без offset пропускается 0 товаров. При limit=0 или отсутствии результатов возвращается пустой массив. У товара без изображения поле image_id равно null.
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param limit query int64 false "Maximum number of products" minimum(0) example(10)
// @Param offset query int64 false "Number of available products to skip" minimum(0) default(0) example(20)
// @Success 200 {array} dto.ProductResponse "Available products, or an empty array"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve products"
// @Router /products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var request dto.ListProductsRequest
	if query.Has("limit") {
		limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/products?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Limit = &limit
	}
	if query.Has("offset") {
		offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/products?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Offset = offset
	}

	err := request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	products, err := h.repository.List(r.Context(), request.Limit, request.Offset)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, mapper.ProductToWeb(product))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при получении списка товаров: %v", err)
	}
}

// @Summary Decrease product stock
// @Description Уменьшает остаток товара на quantity, вплоть до нуля. Дата последней закупки и другие поля не меняются. При недостаточном остатке возвращает 409 без списания.
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Param stock body dto.DecreaseProductStockRequest true "Stock decrease"
// @Success 204 "Product stock decreased"
// @Failure 400 {object} dto.ErrorResponse "Invalid product ID, JSON or quantity"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Product not found"
// @Failure 409 {object} dto.ErrorResponse "Insufficient product stock"
// @Failure 500 {object} dto.ErrorResponse "Failed to decrease product stock"
// @Router /products/{id}/stock [patch]
func (h *ProductHandler) DecreaseStock(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	_, err := uuid.Parse(productID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	var request dto.DecreaseProductStockRequest
	err = decodeJSONBody(r.Body, &request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	found, decreased, err := h.repository.DecreaseStock(r.Context(), productID, request.Quantity)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !found {
		writeError(w, errors.New("товар не найден"), http.StatusNotFound)
		return
	}
	if !decreased {
		writeError(w, errors.New("недостаточный остаток товара"), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
