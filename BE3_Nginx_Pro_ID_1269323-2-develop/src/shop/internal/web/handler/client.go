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

type ClientHandler struct {
	repository domain.ClientRepository
}

func NewClientHandler(repository domain.ClientRepository) *ClientHandler {
	return &ClientHandler{repository: repository}
}

// @Summary Create client
// @Description UUID клиента генерирует сервер. Дату регистрации передаёт вызывающая сторона. Адрес должен уже существовать. Поля client_name, client_surname и gender не могут быть пустыми или состоять только из пробелов.
// @Tags Clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param client body dto.CreateClientRequest true "Client"
// @Success 201 {object} dto.ClientResponse "Client created"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Failed to create client, including a nonexistent address"
// @Router /clients [post]
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateClientRequest
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

	client, err := mapper.ClientToDomain(request)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	err = h.repository.Create(r.Context(), client)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := mapper.ClientToWeb(client)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при создании клиента: %v", err)
	}
}

// @Summary Delete client
// @Description Связанный адрес сохраняется.
// @Tags Clients
// @Security BearerAuth
// @Produce json
// @Param id path string true "Client ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Success 204 "Client deleted"
// @Failure 400 {object} dto.ErrorResponse "Invalid client ID"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Client not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete client"
// @Router /clients/{id} [delete]
func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	clientID := r.PathValue("id")
	_, err := uuid.Parse(clientID)
	if err != nil {
		writeError(w, errors.New("поле id должно содержать корректный UUID"), http.StatusBadRequest)
		return
	}

	deleted, err := h.repository.Delete(r.Context(), clientID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !deleted {
		writeError(w, errors.New("клиент не найден"), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary List or find clients
// @Description Без параметров поиска возвращает клиентов, отсортированных по ID, с необязательной пагинацией. Если передан хотя бы один параметр поиска, оба должны быть непустыми и не состоять только из пробелов. Поиск выполняется по точному совпадению с учётом регистра; пробелы по краям участвуют в поиске. Поиск нельзя совмещать с limit или offset. Без limit количество результатов не ограничено; без offset пропускается 0 клиентов. При limit=0 или отсутствии результатов возвращается пустой массив.
// @Tags Clients
// @Security BearerAuth
// @Produce json
// @Param client_name query string false "Client name; required together with client_surname for search" example(Иван)
// @Param client_surname query string false "Client surname; required together with client_name for search" example(Иванов)
// @Param limit query int64 false "Maximum number of clients; pagination only" minimum(0) example(10)
// @Param offset query int64 false "Number of clients to skip; pagination only" minimum(0) default(0) example(20)
// @Success 200 {array} dto.ClientResponse "Clients, or an empty array"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters or search combined with pagination"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve clients"
// @Router /clients [get]
func (h *ClientHandler) Get(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Has("client_name") || query.Has("client_surname") {
		if query.Has("limit") || query.Has("offset") {
			writeError(w, errors.New("поиск по имени и фамилии нельзя совмещать с пагинацией. Пример запроса: /api/v1/clients?client_name=Иван&client_surname=Иванов"), http.StatusBadRequest)
			return
		}
		h.FindByName(w, r)
		return
	}

	h.List(w, r)
}

func (h *ClientHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	request := dto.FindClientsRequest{
		ClientName:    query.Get("client_name"),
		ClientSurname: query.Get("client_surname"),
	}
	err := request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	clients, err := h.repository.FindByName(r.Context(), request.ClientName, request.ClientSurname)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]dto.ClientResponse, 0, len(clients))
	for _, client := range clients {
		response = append(response, mapper.ClientToWeb(client))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при поиске клиентов: %v", err)
	}
}

func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var request dto.ListClientsRequest
	if query.Has("limit") {
		limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр limit должен быть целым неотрицательным числом. Пример запроса: /api/v1/clients?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Limit = &limit
	}
	if query.Has("offset") {
		offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {
			writeError(w, errors.New("параметр offset должен быть целым неотрицательным числом. Пример запроса: /api/v1/clients?limit=10&offset=20"), http.StatusBadRequest)
			return
		}
		request.Offset = offset
	}

	err := request.Validate()
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	clients, err := h.repository.List(r.Context(), request.Limit, request.Offset)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]dto.ClientResponse, 0, len(clients))
	for _, client := range clients {
		response = append(response, mapper.ClientToWeb(client))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при получении списка клиентов: %v", err)
	}
}

// @Summary Update client address
// @Description При совпадении страны, города и улицы возвращает текущий адрес с прежним UUID. При отличии создаёт новый адрес и привязывает его к клиенту. Старый адрес сохраняется. Страна, город и улица не могут быть пустыми или состоять только из пробелов.
// @Tags Clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Client ID" Format(uuid) example(11111111-1111-4111-8111-111111111111)
// @Param address body dto.CreateAddressRequest true "New address"
// @Success 200 {object} dto.AddressResponse "Current client address"
// @Failure 400 {object} dto.ErrorResponse "Invalid client ID, JSON or required address fields"
// @Failure 401 {object} dto.ErrorResponse "Missing or invalid access token"
// @Failure 404 {object} dto.ErrorResponse "Client not found; no address is created"
// @Failure 500 {object} dto.ErrorResponse "Failed to update client address"
// @Router /clients/{id}/address [put]
func (h *ClientHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	clientID := r.PathValue("id")
	_, err := uuid.Parse(clientID)
	if err != nil {
		writeError(w, errors.New("идентификатор клиента должен содержать корректный UUID"), http.StatusBadRequest)
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
	address, updated, err := h.repository.UpdateAddress(r.Context(), clientID, address)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !updated {
		writeError(w, errors.New("клиент не найден"), http.StatusNotFound)
		return
	}

	response := mapper.AddressToWeb(address)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ошибка записи ответа при смене адреса клиента: %v", err)
	}
}
