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

type SupplierHandler struct {
	repository domain.SupplierRepository
}

func NewSupplierHandler(repository domain.SupplierRepository) *SupplierHandler {
	return &SupplierHandler{repository: repository}
}

// @Summary Create supplier
// @Description UUID поставщика генерирует сервер. Адрес должен уже существовать. Название и телефон не могут состоять только из пробелов; формат телефона не ограничен.
// @Tags Suppliers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param supplier body dto.CreateSupplierRequest true "Supplier"
// @Success 201 {object} dto.SupplierResponse "Supplier created"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Failed to create supplier, including a nonexistent address"
// @Router /suppliers [post]
func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateSupplierRequest
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

	supplier := mapper.SupplierToDomain(request)
	err = h.repository.Create(r.Context(), supplier)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := mapper.SupplierToWeb(supplier)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при создании поставщика: %v", err)
	}
}

// @Summary Delete supplier
// @Description Связанный адрес сохраняется. Наличие любого товара, в том числе с нулевым остатком, препятствует удалению и приводит к ошибке 500.
// @Tags Suppliers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Supplier ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Success 204 "Supplier deleted"
// @Failure 400 {object} dto.ErrorResponse "Invalid supplier ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Supplier not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete supplier, including a supplier with products"
// @Router /suppliers/{id} [delete]
func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	supplierID := r.PathValue("id")
	_, err := uuid.Parse(supplierID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	deleted, err := h.repository.Delete(r.Context(), supplierID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !deleted {
		writeError(w, errors.New("поставщик не найден"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get supplier by ID
// @Tags Suppliers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Supplier ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Success 200 {object} dto.SupplierResponse "Supplier"
// @Failure 400 {object} dto.ErrorResponse "Invalid supplier ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Supplier not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve supplier"
// @Router /suppliers/{id} [get]
func (h *SupplierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	supplierID := r.PathValue("id")
	_, err := uuid.Parse(supplierID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	supplier, found, err := h.repository.GetByID(r.Context(), supplierID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !found {
		writeError(w, errors.New("поставщик не найден"), http.StatusNotFound)
		return
	}

	response := mapper.SupplierToWeb(supplier)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при получении поставщика: %v", err)
	}
}

// @Summary List suppliers
// @Description Возвращает поставщиков, отсортированных по ID по возрастанию, независимо от наличия товаров и их остатка. Параметры пагинации необязательны и могут использоваться отдельно. Без limit количество результатов не ограничено; без offset пропускается 0 поставщиков. При limit=0 или отсутствии результатов возвращается пустой массив.
// @Tags Suppliers
// @Security BearerAuth
// @Produce json
// @Param limit query int64 false "Maximum number of suppliers" minimum(0) example(10)
// @Param offset query int64 false "Number of suppliers to skip" minimum(0) default(0) example(20)
// @Success 200 {array} dto.SupplierResponse "Suppliers, or an empty array"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve suppliers"
// @Router /suppliers [get]
func (h *SupplierHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var request dto.ListSuppliersRequest
	if query.Has("limit") {
		limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/suppliers?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Limit = &limit
	}
	if query.Has("offset") {
		offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/suppliers?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Offset = offset
	}

	err := request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	suppliers, err := h.repository.List(r.Context(), request.Limit, request.Offset)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]dto.SupplierResponse, 0, len(suppliers))
	for _, supplier := range suppliers {
		response = append(response, mapper.SupplierToWeb(supplier))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при получении списка поставщиков: %v", err)
	}
}

// @Summary Update supplier address
// @Description При совпадении страны, города и улицы возвращает текущий адрес с прежним UUID. При отличии создаёт новый адрес и привязывает его к поставщику. Старый адрес сохраняется. Страна, город и улица не могут быть пустыми или состоять только из пробелов.
// @Tags Suppliers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Param address body dto.CreateAddressRequest true "New address"
// @Success 200 {object} dto.AddressResponse "Current supplier address"
// @Failure 400 {object} dto.ErrorResponse "Invalid supplier ID, JSON or required address fields"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Supplier not found; no address is created"
// @Failure 500 {object} dto.ErrorResponse "Failed to update supplier address"
// @Router /suppliers/{id}/address [put]
func (h *SupplierHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	supplierID := r.PathValue("id")
	_, err := uuid.Parse(supplierID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	var request dto.CreateAddressRequest
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

	address := mapper.CreateAddressToDomain(request)
	address, updated, err := h.repository.UpdateAddress(r.Context(), supplierID, address)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !updated {
		writeError(w, errors.New("поставщик не найден"), http.StatusNotFound)
		return
	}

	response := mapper.AddressToWeb(address)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ошибка записи ответа при смене адреса поставщика: %v", err)
	}
}
