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

type AddressHandler struct {
	repository domain.AddressRepository
}

func NewAddressHandler(repository domain.AddressRepository) *AddressHandler {
	return &AddressHandler{repository: repository}
}

// @Summary Получить список адресов
// @Description Возвращает все адреса, отсортированные по ID по возрастанию. Параметры пагинации необязательны и могут использоваться отдельно. Без limit количество результатов не ограничено; без offset пропускается 0 адресов. При limit=0 или отсутствии результатов возвращается пустой массив.
// @Tags Addresses
// @Security BearerAuth
// @Produce json
// @Param limit query int64 false "Максимальное количество адресов" minimum(0) example(10)
// @Param offset query int64 false "Количество адресов, которые нужно пропустить" minimum(0) default(0) example(20)
// @Success 200 {array} dto.AddressResponse "Список адресов или пустой массив"
// @Failure 400 {object} dto.ErrorResponse "Некорректные параметры запроса"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Ошибка получения списка адресов"
// @Router /addresses [get]
func (h *AddressHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var request dto.ListAddressesRequest
	if query.Has("limit") {
		limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/addresses?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Limit = &limit
	}
	if query.Has("offset") {
		offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/addresses?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Offset = offset
	}

	err := request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	addresses, err := h.repository.List(r.Context(), request.Limit, request.Offset)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]dto.AddressResponse, 0, len(addresses))
	for _, address := range addresses {
		response = append(response, mapper.AddressToWeb(address))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при получении списка адресов: %v", err)
	}
}

// @Summary Create address
// @Description Страна, город и улица не могут быть пустыми или состоять только из пробелов.
// @Tags Addresses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param address body dto.CreateAddressRequest true "Address"
// @Success 201 {object} dto.AddressResponse "Address created"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /addresses [post]
func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateAddressRequest
	if err := decodeJSONBody(r.Body, &request); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	address := mapper.CreateAddressToDomain(request)
	if err := h.repository.Create(r.Context(), address); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := mapper.AddressToWeb(address)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ошибка записи ответа при создании адреса: %v", err)
	}
}

// @Summary Update address
// @Description UUID адреса сохраняется. Изменение адреса действует для всех связанных клиентов и поставщиков. Страна, город и улица не могут быть пустыми или состоять только из пробелов.
// @Tags Addresses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Address ID" Format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Param address body dto.ReplaceAddressRequest true "Address"
// @Success 204 "Address updated"
// @Failure 400 {object} dto.ErrorResponse "Invalid address ID, JSON or required fields"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Address not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to update address"
// @Router /addresses/{id} [put]
func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	addressID := r.PathValue("id")
	if _, err := uuid.Parse(addressID); err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	var request dto.ReplaceAddressRequest
	if err := decodeJSONBody(r.Body, &request); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := request.Validate(); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	address := mapper.ReplaceAddressToDomain(addressID, request)
	updated, err := h.repository.Update(r.Context(), address)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !updated {
		writeError(w, errors.New("адрес не найден"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete address
// @Tags Addresses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Address ID" Format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 204 "Address deleted"
// @Failure 400 {object} dto.ErrorResponse "Invalid address ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Address not found"
// @Failure 409 {object} dto.ErrorResponse "Address is in use"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete address"
// @Router /addresses/{id} [delete]
func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	addressID := r.PathValue("id")
	if _, err := uuid.Parse(addressID); err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	deleted, found, err := h.repository.Delete(r.Context(), addressID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !found {
		writeError(w, errors.New("адрес не найден"), http.StatusNotFound)
		return
	}
	if !deleted {
		writeError(w, errors.New("адрес используется"), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
