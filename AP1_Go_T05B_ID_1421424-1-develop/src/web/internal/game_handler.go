package internal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"tictactoe/domain"
)

type GameHandler struct {
	service     domain.GameService
	repository  domain.GameRepository
	userService domain.UserService
}

func NewGameHandler(service domain.GameService, repository domain.GameRepository, userService domain.UserService) *GameHandler {
	result := GameHandler{}

	result.service = service
	result.repository = repository
	result.userService = userService

	return &result
}

func (h *GameHandler) GetAvailableGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.repository.GetAvailableGames(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]GameResponse, 0, len(games))
	for _, game := range games {
		response = append(response, toWeb(game, false))
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи списка доступных игр: %v", err)
	}
}

func (h *GameHandler) GetCompletedGames(w http.ResponseWriter, r *http.Request) {
	value := r.Context().Value(userUUIDKey)
	userUUID, ok := value.(string)
	if userUUID == "" || !ok {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	games, err := h.service.GetCompletedGames(r.Context(), userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]GameResponse, 0, len(games))
	for _, game := range games {
		response = append(response, toWeb(game, true))
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи списка завершенных игр: %v", err)
	}
}

func (h *GameHandler) GetTopPlayers(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil || n <= 0 {
		http.Error(w, "параметр n должен быть целым числом больше нуля", http.StatusBadRequest)
		return
	}

	statistics, err := h.service.GetTopPlayers(r.Context(), n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TopPlayerResponse, 0, len(statistics))
	for _, player := range statistics {
		user, err := h.userService.GetByUUID(r.Context(), player.UserUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response = append(response, TopPlayerResponse{
			UUID:     player.UserUUID,
			Login:    user.Login,
			WinRatio: player.WinRatio,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи списка лучших игроков: %v", err)
	}
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameUUID := r.PathValue("id")
	game, err := h.repository.Get(r.Context(), gameUUID)
	if err != nil {
		if errors.Is(err, domain.ErrGameNotFound) {
			http.Error(w, domain.ErrGameNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gameOver := h.service.IsFinished(r.Context(), game)
	response := toWeb(game, gameOver)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи игры: %v", err)
	}
}

func (h *GameHandler) HandleGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	value := r.Context().Value(userUUIDKey)
	userUUID, ok := value.(string)
	if userUUID == "" || !ok {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

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

	newGame, err = h.service.ApplyHumanMove(r.Context(), savedGame, newGame, userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game := newGame
	if game.PlayerX != "" && game.PlayerX != domain.ComputerUUID &&
		game.PlayerO == domain.ComputerUUID &&
		game.State.Phase == domain.GamePhaseMove &&
		game.State.PlayerUUID == domain.ComputerUUID {
		game, err = h.service.NextMove(r.Context(), game)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = h.repository.Save(r.Context(), game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gameOver := h.service.IsFinished(r.Context(), game)

	response := toWeb(game, gameOver)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа после хода: %v", err)
	}
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	value := r.Context().Value(userUUIDKey)
	userUUID, ok := value.(string)
	if userUUID == "" || !ok {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var gameRequest CreateGameRequest
	err := json.NewDecoder(r.Body).Decode(&gameRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game, err := h.service.CreateGame(r.Context(), domain.GameMode(gameRequest.Mode), userUUID)
	if err != nil {
		if errors.Is(err, domain.ErrIncorrectGameMode) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := toWeb(game, false)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при создании игры: %v", err)
	}
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	value := r.Context().Value(userUUIDKey)
	userUUID, ok := value.(string)
	if userUUID == "" || !ok {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	gameUUID := r.PathValue("id")
	if gameUUID == "" {
		http.Error(w, "передан пустой gameUUID", http.StatusBadRequest)
		return
	}

	game, err := h.service.JoinGame(r.Context(), gameUUID, userUUID)
	if err != nil {
		if errors.Is(err, domain.ErrGameNotFound) {
			http.Error(w, domain.ErrGameNotFound.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrAlreadyInGame) {
			http.Error(w, domain.ErrAlreadyInGame.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, domain.ErrGameUnavailable) {
			http.Error(w, domain.ErrGameUnavailable.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := toWeb(game, false)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ошибка записи ответа при присоединении к игре: %v", err)
	}
}
