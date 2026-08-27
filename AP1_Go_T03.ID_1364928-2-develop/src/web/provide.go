package web

import (
	"net/http"
	"tictactoe/domain"
	"tictactoe/web/internal"
)

type Handler struct {
	handler *internal.GameHandler
}

func NewGameHandler(service domain.GameService, repo domain.GameRepository) *Handler {
	return &Handler{
		handler: internal.NewGameHandler(service, repo),
	}
}

func (h *Handler) HandleGame(w http.ResponseWriter, r *http.Request) {
	h.handler.HandleGame(w, r)
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	h.handler.CreateGame(w, r)
}
