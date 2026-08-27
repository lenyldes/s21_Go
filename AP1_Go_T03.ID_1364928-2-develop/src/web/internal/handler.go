package internal

import (
	"encoding/json"
	"net/http"
	"tictactoe/domain"

	"github.com/google/uuid"
)

type GameHandler struct {
	service    domain.GameService
	repository domain.GameRepository
}

func NewGameHandler(service domain.GameService, repository domain.GameRepository) *GameHandler {
	result := GameHandler{}

	result.service = service
	result.repository = repository

	return &result
}

func (h *GameHandler) HandleGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	UUID := r.PathValue("id")

	var req GameRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	newGame := toDomain(req, UUID)

	savedGame, err := h.repository.Get(r.Context(), UUID)
	if err != nil {
		http.Error(w, "Игра не найдена", http.StatusNotFound)
		return
	}

	err = h.service.Validate(r.Context(), savedGame, newGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.service.IsFinished(r.Context(), newGame) {
		err = h.repository.Save(r.Context(), newGame)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := toWeb(newGame, true)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	nextMove, err := h.service.NextMove(r.Context(), newGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repository.Save(r.Context(), nextMove)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gameOver := h.service.IsFinished(r.Context(), nextMove)

	response := toWeb(nextMove, gameOver)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	newUUID := uuid.New().String()

	newGame := domain.Game{UUID: newUUID, Board: domain.Board{}}

	crGame, err := h.repository.Create(r.Context(), newGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := toWeb(crGame, false)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
